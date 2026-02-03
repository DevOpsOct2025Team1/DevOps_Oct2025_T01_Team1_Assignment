import pytest
from unittest.mock import Mock, patch
import grpc

from file_service.service import FileService, get_user_id
from file.v1 import file_pb2


class TestGetUserId:

    def test_get_user_id_success(self):
        context = Mock()
        context.invocation_metadata.return_value = [("user-id", "test-user-123")]
        
        user_id = get_user_id(context)
        
        assert user_id == "test-user-123"

    def test_get_user_id_missing(self):
        context = Mock()
        context.invocation_metadata.return_value = []
        context.abort.side_effect = Exception("Aborted")

        with pytest.raises(Exception):
            get_user_id(context)

        context.abort.assert_called_once_with(
            grpc.StatusCode.UNAUTHENTICATED,
            "Missing user-id in metadata"
        )


class TestFileService:

    @patch('file_service.service.files_collection')
    def test_create_file(self, mock_collection):
        mock_collection.insert_one.return_value = Mock(inserted_id="507f1f77bcf86cd799439011")
        
        service = FileService()
        context = Mock()
        context.invocation_metadata.return_value = [("user-id", "user-123")]
        
        request = file_pb2.CreateFileRequest(
            filename="test.txt",
            size=1024,
            content_type="text/plain"
        )
        
        response = service.CreateFile(request, context)
        
        assert response.file.id == "507f1f77bcf86cd799439011"
        assert response.file.user_id == "user-123"
        assert response.file.filename == "test.txt"
        assert response.file.size == 1024
        assert response.file.content_type == "text/plain"

    @patch('file_service.service.files_collection')
    def test_list_files(self, mock_collection):
        mock_collection.find.return_value = [
            {
                "_id": "507f1f77bcf86cd799439011",
                "user_id": "user-123",
                "filename": "file1.txt",
                "size": 100,
                "content_type": "text/plain",
                "created_at": 1234567890
            }
        ]
        
        service = FileService()
        context = Mock()
        context.invocation_metadata.return_value = [("user-id", "user-123")]
        
        request = file_pb2.ListFilesRequest()
        response = service.ListFiles(request, context)
        
        assert len(response.files) == 1
        assert response.files[0].filename == "file1.txt"

    @patch('file_service.service.files_collection')
    def test_delete_file(self, mock_collection):
        mock_collection.delete_one.return_value = Mock(deleted_count=1)
        
        service = FileService()
        context = Mock()
        context.invocation_metadata.return_value = [("user-id", "user-123")]
        
        request = file_pb2.DeleteFileRequest(id="507f1f77bcf86cd799439011")
        response = service.DeleteFile(request, context)
        
        assert response.success is True
