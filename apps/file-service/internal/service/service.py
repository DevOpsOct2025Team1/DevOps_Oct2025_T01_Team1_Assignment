import time
from bson import ObjectId
from bson.errors import InvalidId
import grpc

from file.v1 import file_pb2, file_pb2_grpc
from internal.store.db import files_collection

def get_user_id(context):
    #get user id from context metadata
    metadata = dict(context.invocation_metadata())
    user_id = metadata.get("user-id")
    if user_id is None:
        context.abort(grpc.StatusCode.UNAUTHENTICATED, "Missing user-id in metadata")
    return user_id

class FileService(file_pb2_grpc.FileServiceServicer):
    #implementing the File Service


    # Create File
    def CreateFile(self, request, context):
        user_id = get_user_id(context)

        doc = {
            "user_id": user_id,
            "filename": request.filename,
            "size": request.size,
            "content_type": request.content_type,
            "created_at": int(time.time())
        }

        result = files_collection.insert_one(doc)

        return file_pb2.FileResponse(
            file=file_pb2.File(
                id=str(result.inserted_id),
                user_id=user_id,
                filename=doc["filename"],
                size=doc["size"],
                content_type=doc["content_type"],
                created_at=doc["created_at"]
            ) #return success message
        )

    #get list of files for user
    def ListFiles(self, request, context):
        user_id = get_user_id(context)

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

    #get a specific file for user by id
    def GetFile(self, request, context):
        user_id = get_user_id(context)

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

    #delete a specific file for user by id
    def DeleteFile(self, request, context):
        user_id = get_user_id(context)

        try:
            res = files_collection.delete_one({"_id": ObjectId(request.id), "user_id": user_id})
        except InvalidId:
            context.abort(grpc.StatusCode.INVALID_ARGUMENT, "Invalid file id format")
        return file_pb2.DeleteFileResponse(success=res.deleted_count == 1)