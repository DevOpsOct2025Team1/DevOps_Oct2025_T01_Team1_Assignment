from fastapi import FastAPI, UploadFile, File, Depends, HTTPException, Header, Response
from fastapi.responses import StreamingResponse, JSONResponse
from fastapi.middleware.cors import CORSMiddleware
from typing import List, Optional
import uvicorn
import grpc
from contextlib import asynccontextmanager

from file_service.service import FileService
from file_service.auth_client import AuthClient
from file_service.store import files_collection, s3_client, generate_s3_key
from file_service.config import S3_BUCKET_NAME, HTTP_PORT
from bson import ObjectId
from botocore.exceptions import ClientError

# Global auth client
auth_client = None

@asynccontextmanager
async def lifespan(app: FastAPI):
    global auth_client
    auth_client = AuthClient()
    yield
    if auth_client:
        auth_client.close()

app = FastAPI(lifespan=lifespan)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Adjust for production
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Dependency to get current user
async def get_current_user(authorization: str = Header(None)):
    if not authorization:
        raise HTTPException(status_code=401, detail="Missing Authorization header")
    
    scheme, _, token = authorization.partition(" ")
    if scheme.lower() != "bearer" or not token:
         raise HTTPException(status_code=401, detail="Invalid Authorization header format")

    response = None
    if auth_client:
        response = auth_client.validate_token(token)
    
    if not response or not response.valid:
        raise HTTPException(status_code=401, detail="Invalid token")
    
    return response.user_id

@app.post("/api/files")
async def upload_file(file: UploadFile = File(...), user_id: str = Depends(get_current_user)):
    try:
        file_id = str(ObjectId())
        filename = file.filename
        content_type = file.content_type
        s3_key = generate_s3_key(user_id, file_id, filename)
        
        # Upload to S3
        s3_client.upload_fileobj(
            file.file,
            S3_BUCKET_NAME,
            s3_key,
            ExtraArgs={'ContentType': content_type}
        )
        
        # Save metadata to MongoDB
        file_size = file.size # Note: spooled file size might need checking
        # If file.size is not available or reliable (e.g. chunked), we might need to seek/tell
        # But for UploadFile it usually works if spooled.
        
        doc = {
            "_id": ObjectId(file_id),
            "user_id": user_id,
            "filename": filename,
            "size": file_size,
            "content_type": content_type,
            "s3_key": s3_key,
            "created_at": int(file_id[:8], 16) # timestamp from ObjectId
        }
        files_collection.insert_one(doc)
        
        return {
            "file": {
                "id": file_id,
                "filename": filename,
                "size": file_size,
                "content_type": content_type,
                "created_at": doc["created_at"]
            }
        }
    except Exception as e:
        print(f"Upload failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/api/files")
async def list_files(user_id: str = Depends(get_current_user)):
    cursor = files_collection.find({"user_id": user_id})
    files = []
    for doc in cursor:
        files.append({
            "id": str(doc["_id"]),
            "filename": doc["filename"],
            "size": doc["size"],
            "content_type": doc["content_type"],
            "created_at": doc.get("created_at")
        })
    return {"files": files}

@app.get("/api/files/{file_id}/download")
async def download_file(file_id: str, user_id: str = Depends(get_current_user)):
    try:
        doc = files_collection.find_one({"_id": ObjectId(file_id), "user_id": user_id})
    except:
         raise HTTPException(status_code=400, detail="Invalid file ID")

    if not doc:
        raise HTTPException(status_code=404, detail="File not found")
        
    s3_key = doc.get("s3_key")
    if not s3_key:
         raise HTTPException(status_code=500, detail="File content missing")

    try:
        # Get stream from S3
        response = s3_client.get_object(Bucket=S3_BUCKET_NAME, Key=s3_key)
        
        def iterfile():
            yield from response['Body']

        return StreamingResponse(
            iterfile(),
            media_type=doc["content_type"],
            headers={"Content-Disposition": f'attachment; filename="{doc["filename"]}"'}
        )
    except ClientError as e:
        raise HTTPException(status_code=500, detail=f"S3 Error: {str(e)}")

@app.delete("/api/files/{file_id}")
async def delete_file(file_id: str, user_id: str = Depends(get_current_user)):
    try:
        doc = files_collection.find_one({"_id": ObjectId(file_id), "user_id": user_id})
    except:
         raise HTTPException(status_code=400, detail="Invalid file ID")
         
    if not doc:
        return JSONResponse(status_code=404, content={"error": "File not found"})

    # Delete from S3
    s3_key = doc.get("s3_key")
    if s3_key:
        try:
            s3_client.delete_object(Bucket=S3_BUCKET_NAME, Key=s3_key)
        except ClientError:
            pass # Continue to delete metadata
            
    files_collection.delete_one({"_id": ObjectId(file_id)})
    return {"success": True}

def run_http_server():
    print(f"Starting HTTP server on port {HTTP_PORT}")
    uvicorn.run(app, host="0.0.0.0", port=int(HTTP_PORT))
