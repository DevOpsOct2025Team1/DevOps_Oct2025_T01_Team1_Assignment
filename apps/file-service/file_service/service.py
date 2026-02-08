import time
import io
from bson import ObjectId
from bson.errors import InvalidId
import grpc
from botocore.exceptions import ClientError
from file.v1 import file_pb2, file_pb2_grpc
from file_service.store import files_collection, s3_client, generate_s3_key
from file_service.config import S3_BUCKET_NAME
from file_service.auth_client import AuthClient


def get_user_id(context, auth_client):
    """Extract and validate authorization token from gRPC metadata."""
    metadata = dict(context.invocation_metadata())
    authorization = metadata.get("authorization")

    if not authorization:
        context.abort(grpc.StatusCode.UNAUTHENTICATED, "Missing authorization header")

    parts = authorization.split(" ")
    if len(parts) != 2 or parts[0].lower() != "bearer":
        context.abort(grpc.StatusCode.UNAUTHENTICATED, "Invalid authorization header format")

    token = parts[1]

    response = auth_client.validate_token(token)
    if not response or not response.valid:
        context.abort(grpc.StatusCode.UNAUTHENTICATED, "Invalid token")

    if not response.user:
        context.abort(grpc.StatusCode.UNAUTHENTICATED, "Invalid token: no user")

    return response.user.id


class FileService(file_pb2_grpc.FileServiceServicer):
    def __init__(self, auth_client: AuthClient):
        self.auth_client = auth_client

    def CreateFile(self, request, context):
        user_id = get_user_id(context, self.auth_client)

        filename = request.filename
        size = request.size
        content_type = request.content_type or "application/octet-stream"

        file_id = ObjectId()
        s3_key = generate_s3_key(user_id, str(file_id), filename)

        doc = {
            "_id": file_id,
            "user_id": user_id,
            "filename": filename,
            "size": size,
            "content_type": content_type,
            "s3_key": s3_key,
            "created_at": int(time.time()),
        }
        result = files_collection.insert_one(doc)

        return file_pb2.FileResponse(
            file=file_pb2.File(
                id=str(result.inserted_id),
                user_id=user_id,
                filename=filename,
                size=size,
                content_type=content_type,
                created_at=doc["created_at"],
            )
        )

    def UploadFile(self, request_iterator, context):
        """Stream upload file to S3 and save metadata."""
        user_id = get_user_id(context, self.auth_client)

        # First message should contain metadata
        MAX_FILES_PER_USER = 20
        MAX_FILE_SIZE = 2 * 1024 * 1024 * 1024

        file_count = files_collection.count_documents({"user_id": user_id})
        if file_count >= MAX_FILES_PER_USER:
            context.abort(
                grpc.StatusCode.RESOURCE_EXHAUSTED,
                "Maximum file limit reached (20 files per user)"
            )

        try:
            first_request = next(request_iterator)
            if not first_request.HasField("metadata"):
                context.abort(grpc.StatusCode.INVALID_ARGUMENT, "First message must contain metadata")
            
            metadata = first_request.metadata
            filename = metadata.filename
            content_type = metadata.content_type or "application/octet-stream"
            
            # Generate unique file ID and S3 key
            file_id = str(ObjectId())
            s3_key = generate_s3_key(user_id, file_id, filename)
            
            # Collect file chunks
            file_buffer = io.BytesIO()
            total_size = 0
            for request in request_iterator:
                if request.HasField("chunk"):
                    chunk_size = len(request.chunk)
                    total_size += chunk_size

                    if total_size > MAX_FILE_SIZE:
                        context.abort(
                            grpc.StatusCode.RESOURCE_EXHAUSTED,
                            "File size exceeds maximum allowed size (2GB)"
                        )

                    file_buffer.write(request.chunk)
            
            # Upload to S3
            file_buffer.seek(0)
            file_size = file_buffer.getbuffer().nbytes
            
            try:
                s3_client.upload_fileobj(
                    file_buffer,
                    S3_BUCKET_NAME,
                    s3_key,
                    ExtraArgs={'ContentType': content_type}
                )
            except ClientError as e:
                context.abort(grpc.StatusCode.INTERNAL, f"Failed to upload file to S3: {str(e)}")
            
            # Save metadata to MongoDB
            doc = {
                "_id": ObjectId(file_id),
                "user_id": user_id,
                "filename": filename,
                "size": file_size,
                "content_type": content_type,
                "s3_key": s3_key,
                "created_at": int(time.time())
            }
            files_collection.insert_one(doc)
            
            return file_pb2.FileResponse(
                file=file_pb2.File(
                    id=file_id,
                    user_id=user_id,
                    filename=filename,
                    size=file_size,
                    content_type=content_type,
                    created_at=doc["created_at"]
                )
            )
        except StopIteration:
            context.abort(grpc.StatusCode.INVALID_ARGUMENT, "Empty upload stream")

    def ListFiles(self, request, context):
        user_id = get_user_id(context, self.auth_client)

        docs = files_collection.find({"user_id": user_id})

        return file_pb2.ListFilesResponse(
            files=[
                file_pb2.File(
                    id=str(d["_id"]),
                    user_id=d["user_id"],
                    filename=d["filename"],
                    size=d["size"],
                    content_type=d["content_type"],
                    created_at=d["created_at"]
                ) for d in docs
            ]
        )

    def GetFile(self, request, context):
        user_id = get_user_id(context, self.auth_client)

        try:
            doc = files_collection.find_one({"_id": ObjectId(request.id), "user_id": user_id})
        except InvalidId:
            context.abort(grpc.StatusCode.INVALID_ARGUMENT, "Invalid file id format")

        if not doc:
            context.abort(grpc.StatusCode.NOT_FOUND, "File not found")

        return file_pb2.FileResponse(
            file=file_pb2.File(
                id=str(doc["_id"]),
                user_id=doc["user_id"],
                filename=doc["filename"],
                size=doc["size"],
                content_type=doc["content_type"],
                created_at=doc["created_at"]
            )
        )


    def DownloadFile(self, request, context):
        """Stream download file from S3."""
        user_id = get_user_id(context, self.auth_client)
        
        try:
            doc = files_collection.find_one({"_id": ObjectId(request.id), "user_id": user_id})
        except InvalidId:
            context.abort(grpc.StatusCode.INVALID_ARGUMENT, "Invalid file id format")
        
        if not doc:
            context.abort(grpc.StatusCode.NOT_FOUND, "File not found")
        
        s3_key = doc.get("s3_key")
        if not s3_key:
            context.abort(grpc.StatusCode.INTERNAL, "File metadata missing S3 key")
        
        try:
            # First, yield metadata
            yield file_pb2.DownloadFileResponse(
                metadata=file_pb2.DownloadFileMetadata(
                    filename=doc["filename"],
                    content_type=doc["content_type"],
                    size=doc["size"]
                )
            )
            
            # Stream file from S3 in chunks
            response = s3_client.get_object(Bucket=S3_BUCKET_NAME, Key=s3_key)
            chunk_size = 64 * 1024  # 64KB chunks
            
            try:
                while True:
                    chunk = response['Body'].read(chunk_size)
                    if not chunk:
                        break
                    yield file_pb2.DownloadFileResponse(chunk=chunk)
            finally:
                response['Body'].close()
                
        except ClientError as e:
            context.abort(grpc.StatusCode.INTERNAL, f"Failed to download file from S3: {str(e)}")

    def DeleteFile(self, request, context):
        user_id = get_user_id(context, self.auth_client)

        try:
            doc = files_collection.find_one({"_id": ObjectId(request.id), "user_id": user_id})
        except InvalidId:
            context.abort(grpc.StatusCode.INVALID_ARGUMENT, "Invalid file id format")
        
        if not doc:
            return file_pb2.DeleteFileResponse(success=False)
        
        # Delete from S3 if s3_key exists
        s3_key = doc.get("s3_key")
        if s3_key:
            try:
                s3_client.delete_object(Bucket=S3_BUCKET_NAME, Key=s3_key)
            except ClientError as e:
                # Log error but continue with MongoDB deletion
                print(f"Warning: Failed to delete file from S3: {str(e)}")
        
        # Delete from MongoDB
        res = files_collection.delete_one({"_id": ObjectId(request.id), "user_id": user_id})
        return file_pb2.DeleteFileResponse(success=res.deleted_count == 1)