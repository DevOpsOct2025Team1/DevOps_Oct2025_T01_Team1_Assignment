import pytest
from unittest.mock import Mock, MagicMock, patch
import grpc

from file_service.service import FileService, get_user_id
from file.v1 import file_pb2
from auth.v1 import auth_pb2
from user.v1 import user_pb2


class TestGetUserId:

    def test_get_user_id_success(self):
        context = Mock()
        context.invocation_metadata.return_value = [("authorization", "Bearer valid-token")]

        auth_client = Mock()
        auth_client.validate_token.return_value = auth_pb2.ValidateTokenResponse(
            valid=True,
            user=user_pb2.User(id="test-user-123", username="testuser", role=user_pb2.Role.ROLE_USER)
        )

        user_id = get_user_id(context, auth_client)

        assert user_id == "test-user-123"
        auth_client.validate_token.assert_called_once_with("valid-token")

    def test_get_user_id_missing_authorization(self):
        context = Mock()
        context.invocation_metadata.return_value = []
        context.abort.side_effect = Exception("Aborted")

        auth_client = Mock()

        with pytest.raises(Exception):
            get_user_id(context, auth_client)

        context.abort.assert_called_once_with(
            grpc.StatusCode.UNAUTHENTICATED,
            "Missing authorization header"
        )

    def test_get_user_id_invalid_format(self):
        context = Mock()
        context.invocation_metadata.return_value = [("authorization", "InvalidToken")]
        context.abort.side_effect = Exception("Aborted")

        auth_client = Mock()

        with pytest.raises(Exception):
            get_user_id(context, auth_client)

        context.abort.assert_called_once_with(
            grpc.StatusCode.UNAUTHENTICATED,
            "Invalid authorization header format"
        )

    def test_get_user_id_auth_service_unavailable(self):
        context = Mock()
        context.invocation_metadata.return_value = [("authorization", "Bearer valid-token")]
        context.abort.side_effect = Exception("Aborted")

        auth_client = Mock()
        auth_client.validate_token.side_effect = grpc.RpcError()

        with pytest.raises(Exception):
            get_user_id(context, auth_client)

        context.abort.assert_called_once()
        assert context.abort.call_args[0][0] == grpc.StatusCode.UNAVAILABLE

    def test_get_user_id_invalid_token(self):
        context = Mock()
        context.invocation_metadata.return_value = [("authorization", "Bearer invalid-token")]
        context.abort.side_effect = Exception("Aborted")

        auth_client = Mock()
        auth_client.validate_token.return_value = auth_pb2.ValidateTokenResponse(
            valid=False
        )

        with pytest.raises(Exception):
            get_user_id(context, auth_client)

        context.abort.assert_called_once_with(
            grpc.StatusCode.UNAUTHENTICATED,
            "Invalid token"
        )


class TestFileService:

    @patch('file_service.service.files_collection')
    def test_create_file(self, mock_collection):
        mock_collection.insert_one.return_value = Mock(inserted_id="507f1f77bcf86cd799439011")

        auth_client = Mock()
        auth_client.validate_token.return_value = auth_pb2.ValidateTokenResponse(
            valid=True,
            user=user_pb2.User(id="user-123", username="testuser", role=user_pb2.Role.ROLE_USER)
        )

        service = FileService(auth_client)
        context = Mock()
        context.invocation_metadata.return_value = [("authorization", "Bearer valid-token")]

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

        auth_client = Mock()
        auth_client.validate_token.return_value = auth_pb2.ValidateTokenResponse(
            valid=True,
            user=user_pb2.User(id="user-123", username="testuser", role=user_pb2.Role.ROLE_USER)
        )

        service = FileService(auth_client)
        context = Mock()
        context.invocation_metadata.return_value = [("authorization", "Bearer valid-token")]
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

        auth_client = Mock()
        auth_client.validate_token.return_value = auth_pb2.ValidateTokenResponse(
            valid=True,
            user=user_pb2.User(id="user-123", username="testuser", role=user_pb2.Role.ROLE_USER)
        )

        service = FileService(auth_client)
        context = Mock()
        context.invocation_metadata.return_value = [("authorization", "Bearer valid-token")]

        request = file_pb2.ListFilesRequest()
        response = service.ListFiles(request, context)

        assert len(response.files) == 1
        assert response.files[0].filename == "file1.txt"

        mock_collection.find.assert_called_once_with({"user_id": "user-123"})

    @patch('file_service.service.files_collection')
    def test_list_files_ignores_user_id_in_request(self, mock_collection):
        mock_collection.find.return_value = []

        auth_client = Mock()
        auth_client.validate_token.return_value = auth_pb2.ValidateTokenResponse(
            valid=True,
            user=user_pb2.User(id="user-123", username="testuser", role=user_pb2.Role.ROLE_USER)
        )

        service = FileService(auth_client)
        context = Mock()
        context.invocation_metadata.return_value = [("authorization", "Bearer valid-token")]

        request = file_pb2.ListFilesRequest()
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

        auth_client = Mock()
        auth_client.validate_token.return_value = auth_pb2.ValidateTokenResponse(
            valid=True,
            user=user_pb2.User(id="user-123", username="testuser", role=user_pb2.Role.ROLE_USER)
        )

        service = FileService(auth_client)
        context = Mock()
        context.invocation_metadata.return_value = [("authorization", "Bearer valid-token")]

        request = file_pb2.GetFileRequest(id="507f1f77bcf86cd799439011")
        response = service.GetFile(request, context)

        assert response.file.filename == "file1.txt"
        assert response.file.user_id == "user-123"

    @patch('file_service.service.files_collection')
    def test_get_file_invalid_id_aborts(self, mock_collection):
        auth_client = Mock()
        auth_client.validate_token.return_value = auth_pb2.ValidateTokenResponse(
            valid=True,
            user=user_pb2.User(id="user-123", username="testuser", role=user_pb2.Role.ROLE_USER)
        )

        service = FileService(auth_client)
        context = Mock()
        context.invocation_metadata.return_value = [("authorization", "Bearer valid-token")]
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

        auth_client = Mock()
        auth_client.validate_token.return_value = auth_pb2.ValidateTokenResponse(
            valid=True,
            user=user_pb2.User(id="user-123", username="testuser", role=user_pb2.Role.ROLE_USER)
        )

        service = FileService(auth_client)
        context = Mock()
        context.invocation_metadata.return_value = [("authorization", "Bearer valid-token")]
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
        mock_collection.find_one.return_value = {"_id": "507f1f77bcf86cd799439011", "user_id": "user-123"}
        mock_collection.delete_one.return_value = Mock(deleted_count=1)

        auth_client = Mock()
        auth_client.validate_token.return_value = auth_pb2.ValidateTokenResponse(
            valid=True,
            user=user_pb2.User(id="user-123", username="testuser", role=user_pb2.Role.ROLE_USER)
        )

        service = FileService(auth_client)
        context = Mock()
        context.invocation_metadata.return_value = [("authorization", "Bearer valid-token")]

        request = file_pb2.DeleteFileRequest(id="507f1f77bcf86cd799439011")
        response = service.DeleteFile(request, context)

        assert response.success is True

    @patch('file_service.service.files_collection')
    def test_delete_file_not_found(self, mock_collection):
        mock_collection.find_one.return_value = None

        auth_client = Mock()
        auth_client.validate_token.return_value = auth_pb2.ValidateTokenResponse(
            valid=True,
            user=user_pb2.User(id="user-123", username="testuser", role=user_pb2.Role.ROLE_USER)
        )

        service = FileService(auth_client)
        context = Mock()
        context.invocation_metadata.return_value = [("authorization", "Bearer valid-token")]
        context.abort.side_effect = Exception("Aborted")

        request = file_pb2.DeleteFileRequest(id="507f1f77bcf86cd799439011")
        with pytest.raises(Exception):
            service.DeleteFile(request, context)

        context.abort.assert_called_once_with(
            grpc.StatusCode.NOT_FOUND,
            "File not found",
        )

    @patch('file_service.service.files_collection')
    def test_delete_file_invalid_id_aborts(self, mock_collection):
        auth_client = Mock()
        auth_client.validate_token.return_value = auth_pb2.ValidateTokenResponse(
            valid=True,
            user=user_pb2.User(id="user-123", username="testuser", role=user_pb2.Role.ROLE_USER)
        )

        service = FileService(auth_client)
        context = Mock()
        context.invocation_metadata.return_value = [("authorization", "Bearer valid-token")]
        context.abort.side_effect = Exception("Aborted")

        request = file_pb2.DeleteFileRequest(id="not-a-valid-objectid")
        with pytest.raises(Exception):
            service.DeleteFile(request, context)

        context.abort.assert_called_once_with(
            grpc.StatusCode.INVALID_ARGUMENT,
            "Invalid file id format",
        )
        mock_collection.find_one.assert_not_called()

class TestBusinessRules:

    @patch('file_service.service.files_collection')
    def test_upload_file_exceeds_max_files_per_user(self, mock_collection):
        mock_collection.count_documents.return_value = 20

        auth_client = Mock()
        auth_client.validate_token.return_value = auth_pb2.ValidateTokenResponse(
            valid=True,
            user=user_pb2.User(id="user-123", username="testuser", role=user_pb2.Role.ROLE_USER)
        )

        service = FileService(auth_client)
        context = Mock()
        context.invocation_metadata.return_value = [("authorization", "Bearer valid-token")]
        context.abort.side_effect = Exception("Aborted")

        def request_iterator():
            yield file_pb2.UploadFileRequest(
                metadata=file_pb2.UploadFileMetadata(
                    filename="test.txt",
                    content_type="text/plain"
                )
            )

        with pytest.raises(Exception):
            service.UploadFile(request_iterator(), context)

        context.abort.assert_called_once_with(
            grpc.StatusCode.RESOURCE_EXHAUSTED,
            "Maximum file limit reached (20 files per user)"
        )

    @pytest.mark.integration
    @patch('file_service.service.files_collection')
    def test_upload_file_exceeds_max_file_size(self, mock_collection):
        mock_collection.count_documents.return_value = 5

        auth_client = Mock()
        auth_client.validate_token.return_value = auth_pb2.ValidateTokenResponse(
            valid=True,
            user=user_pb2.User(id="user-123", username="testuser", role=user_pb2.Role.ROLE_USER)
        )

        service = FileService(auth_client)
        context = Mock()
        context.invocation_metadata.return_value = [("authorization", "Bearer valid-token")]
        context.abort.side_effect = Exception("Aborted")

        MAX_FILE_SIZE = 2 * 1024 * 1024 * 1024

        def request_iterator():
            yield file_pb2.UploadFileRequest(
                metadata=file_pb2.UploadFileMetadata(
                    filename="large.bin",
                    content_type="application/octet-stream"
                )
            )
            chunk_size = 64 * 1024
            total_sent = 0
            while total_sent <= MAX_FILE_SIZE:
                yield file_pb2.UploadFileRequest(chunk=b'x' * chunk_size)
                total_sent += chunk_size

        with pytest.raises(Exception):
            service.UploadFile(request_iterator(), context)

        context.abort.assert_called_with(
            grpc.StatusCode.RESOURCE_EXHAUSTED,
            "File size exceeds maximum allowed size (2GB)"
        )


class TestDownloadFile:

    @patch('file_service.service.s3_client')
    @patch('file_service.service.files_collection')
    def test_download_file_success(self, mock_collection, mock_s3):
        mock_collection.find_one.return_value = {
            "_id": "507f1f77bcf86cd799439011",
            "user_id": "user-123",
            "filename": "test.txt",
            "size": 11,
            "content_type": "text/plain",
            "s3_key": "user-123/507f1f77bcf86cd799439011/test.txt",
            "created_at": 1234567890,
        }

        body_mock = MagicMock()
        body_mock.read.side_effect = [b"hello world", b""]
        mock_s3.get_object.return_value = {"Body": body_mock}

        auth_client = Mock()
        auth_client.validate_token.return_value = auth_pb2.ValidateTokenResponse(
            valid=True,
            user=user_pb2.User(id="user-123", username="testuser", role=user_pb2.Role.ROLE_USER)
        )

        service = FileService(auth_client)
        context = Mock()
        context.invocation_metadata.return_value = [("authorization", "Bearer valid-token")]

        request = file_pb2.DownloadFileRequest(id="507f1f77bcf86cd799439011")
        responses = list(service.DownloadFile(request, context))

        # First response should be metadata
        assert responses[0].HasField("metadata")
        assert responses[0].metadata.filename == "test.txt"
        assert responses[0].metadata.content_type == "text/plain"
        assert responses[0].metadata.size == 11

        # Second response should be the chunk
        assert responses[1].chunk == b"hello world"

        body_mock.close.assert_called_once()

    @patch('file_service.service.s3_client')
    @patch('file_service.service.files_collection')
    def test_download_file_streams_multiple_chunks(self, mock_collection, mock_s3):
        mock_collection.find_one.return_value = {
            "_id": "507f1f77bcf86cd799439011",
            "user_id": "user-123",
            "filename": "big.bin",
            "size": 200,
            "content_type": "application/octet-stream",
            "s3_key": "user-123/507f1f77bcf86cd799439011/big.bin",
            "created_at": 1234567890,
        }

        body_mock = MagicMock()
        body_mock.read.side_effect = [b"a" * 100, b"b" * 100, b""]
        mock_s3.get_object.return_value = {"Body": body_mock}

        auth_client = Mock()
        auth_client.validate_token.return_value = auth_pb2.ValidateTokenResponse(
            valid=True,
            user=user_pb2.User(id="user-123", username="testuser", role=user_pb2.Role.ROLE_USER)
        )

        service = FileService(auth_client)
        context = Mock()
        context.invocation_metadata.return_value = [("authorization", "Bearer valid-token")]

        request = file_pb2.DownloadFileRequest(id="507f1f77bcf86cd799439011")
        responses = list(service.DownloadFile(request, context))

        assert len(responses) == 3  # 1 metadata + 2 chunks
        assert responses[1].chunk == b"a" * 100
        assert responses[2].chunk == b"b" * 100

    @patch('file_service.service.files_collection')
    def test_download_file_not_found(self, mock_collection):
        mock_collection.find_one.return_value = None

        auth_client = Mock()
        auth_client.validate_token.return_value = auth_pb2.ValidateTokenResponse(
            valid=True,
            user=user_pb2.User(id="user-123", username="testuser", role=user_pb2.Role.ROLE_USER)
        )

        service = FileService(auth_client)
        context = Mock()
        context.invocation_metadata.return_value = [("authorization", "Bearer valid-token")]
        context.abort.side_effect = Exception("Aborted")

        request = file_pb2.DownloadFileRequest(id="507f1f77bcf86cd799439011")
        with pytest.raises(Exception):
            list(service.DownloadFile(request, context))

        context.abort.assert_called_once_with(
            grpc.StatusCode.NOT_FOUND,
            "File not found",
        )

    @patch('file_service.service.s3_client')
    @patch('file_service.service.files_collection')
    def test_download_file_s3_error(self, mock_collection, mock_s3):
        from botocore.exceptions import ClientError

        mock_collection.find_one.return_value = {
            "_id": "507f1f77bcf86cd799439011",
            "user_id": "user-123",
            "filename": "test.txt",
            "size": 11,
            "content_type": "text/plain",
            "s3_key": "user-123/507f1f77bcf86cd799439011/test.txt",
            "created_at": 1234567890,
        }

        mock_s3.get_object.side_effect = ClientError(
            {"Error": {"Code": "NoSuchKey", "Message": "Not found"}},
            "GetObject"
        )

        auth_client = Mock()
        auth_client.validate_token.return_value = auth_pb2.ValidateTokenResponse(
            valid=True,
            user=user_pb2.User(id="user-123", username="testuser", role=user_pb2.Role.ROLE_USER)
        )

        service = FileService(auth_client)
        context = Mock()
        context.invocation_metadata.return_value = [("authorization", "Bearer valid-token")]
        context.abort.side_effect = Exception("Aborted")

        request = file_pb2.DownloadFileRequest(id="507f1f77bcf86cd799439011")
        with pytest.raises(Exception):
            list(service.DownloadFile(request, context))

        context.abort.assert_called_once()
        assert context.abort.call_args[0][0] == grpc.StatusCode.INTERNAL
