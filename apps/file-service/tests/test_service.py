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

        mock_collection.insert_one.assert_called_once()
        inserted_doc = mock_collection.insert_one.call_args.args[0]
        assert inserted_doc["user_id"] == "user-123"
        assert inserted_doc["filename"] == "test.txt"
        assert inserted_doc["size"] == 1024
        assert inserted_doc["content_type"] == "text/plain"
        assert isinstance(inserted_doc["created_at"], int)

    @patch('file_service.service.time.time', return_value=1700000000.9)
    @patch('file_service.service.files_collection')
    def test_create_file_sets_created_at_from_time(self, mock_collection, _mock_time):
        mock_collection.insert_one.return_value = Mock(inserted_id="507f1f77bcf86cd799439011")

        service = FileService()
        context = Mock()
        context.invocation_metadata.return_value = [("user-id", "user-123")]
        request = file_pb2.CreateFileRequest(filename="x", size=1, content_type="text/plain")

        response = service.CreateFile(request, context)
        assert response.file.created_at == 1700000000

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

        mock_collection.find.assert_called_once_with({"user_id": "user-123"})

    @patch('file_service.service.files_collection')
    def test_list_files_ignores_user_id_in_request(self, mock_collection):
        mock_collection.find.return_value = []

        service = FileService()
        context = Mock()
        context.invocation_metadata.return_value = [("user-id", "user-123")]

        request = file_pb2.ListFilesRequest(user_id="other-user")
        _ = service.ListFiles(request, context)

        mock_collection.find.assert_called_once_with({"user_id": "user-123"})

    @patch('file_service.service.files_collection')
    def test_get_file_success(self, mock_collection):
        mock_collection.find_one.return_value = {
            "_id": "507f1f77bcf86cd799439011",
            "user_id": "user-123",
            "filename": "file1.txt",
            "size": 100,
            "content_type": "text/plain",
            "created_at": 1234567890,
        }

        service = FileService()
        context = Mock()
        context.invocation_metadata.return_value = [("user-id", "user-123")]

        request = file_pb2.GetFileRequest(id="507f1f77bcf86cd799439011")
        response = service.GetFile(request, context)

        assert response.file.filename == "file1.txt"
        assert response.file.user_id == "user-123"

    @patch('file_service.service.files_collection')
    def test_get_file_invalid_id_aborts(self, mock_collection):
        service = FileService()
        context = Mock()
        context.invocation_metadata.return_value = [("user-id", "user-123")]
        context.abort.side_effect = Exception("Aborted")

        request = file_pb2.GetFileRequest(id="not-a-valid-objectid")
        with pytest.raises(Exception):
            service.GetFile(request, context)

        context.abort.assert_called_once_with(
            grpc.StatusCode.INVALID_ARGUMENT,
            "Invalid file id format",
        )
        mock_collection.find_one.assert_not_called()

    @patch('file_service.service.files_collection')
    def test_get_file_not_found_aborts(self, mock_collection):
        mock_collection.find_one.return_value = None

        service = FileService()
        context = Mock()
        context.invocation_metadata.return_value = [("user-id", "user-123")]
        context.abort.side_effect = Exception("Aborted")

        request = file_pb2.GetFileRequest(id="507f1f77bcf86cd799439011")
        with pytest.raises(Exception):
            service.GetFile(request, context)

        context.abort.assert_called_once_with(
            grpc.StatusCode.NOT_FOUND,
            "File not found",
        )

    @patch('file_service.service.files_collection')
    def test_delete_file(self, mock_collection):
        mock_collection.delete_one.return_value = Mock(deleted_count=1)
        
        service = FileService()
        context = Mock()
        context.invocation_metadata.return_value = [("user-id", "user-123")]
        
        request = file_pb2.DeleteFileRequest(id="507f1f77bcf86cd799439011")
        response = service.DeleteFile(request, context)
        
        assert response.success is True

    @patch('file_service.service.files_collection')
    def test_delete_file_not_found(self, mock_collection):
        mock_collection.delete_one.return_value = Mock(deleted_count=0)

        service = FileService()
        context = Mock()
        context.invocation_metadata.return_value = [("user-id", "user-123")]

        request = file_pb2.DeleteFileRequest(id="507f1f77bcf86cd799439011")
        response = service.DeleteFile(request, context)

        assert response.success is False

    @patch('file_service.service.files_collection')
    def test_delete_file_invalid_id_aborts(self, mock_collection):
        service = FileService()
        context = Mock()
        context.invocation_metadata.return_value = [("user-id", "user-123")]
        context.abort.side_effect = Exception("Aborted")

        request = file_pb2.DeleteFileRequest(id="not-a-valid-objectid")
        with pytest.raises(Exception):
            service.DeleteFile(request, context)

        context.abort.assert_called_once_with(
            grpc.StatusCode.INVALID_ARGUMENT,
            "Invalid file id format",
        )
        mock_collection.delete_one.assert_not_called()
