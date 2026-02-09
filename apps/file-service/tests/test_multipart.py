import pytest
from unittest.mock import Mock, patch
from bson import ObjectId
from file.v1 import file_pb2
from file_service.service import FileService

pytestmark = pytest.mark.integration

@pytest.fixture
def mock_auth_client():
    client = Mock()
    validate_response = Mock()
    validate_response.valid = True
    user = Mock()
    user.id = "test_user_123"
    validate_response.user = user
    client.validate_token.return_value = validate_response
    return client

@pytest.fixture
def mock_context():
    context = Mock()
    context.invocation_metadata.return_value = [
        ("authorization", "Bearer test_token")
    ]
    return context

@pytest.fixture
def file_service(mock_auth_client):
    return FileService(mock_auth_client)

def test_initiate_multipart_upload(file_service, mock_context):
    request = file_pb2.InitiateMultipartUploadRequest(
        filename="large_file.mp4",
        content_type="video/mp4",
        total_size=100 * 1024 * 1024
    )

    with patch('file_service.service.s3_client') as mock_s3, \
         patch('file_service.service.upload_sessions_collection') as mock_sessions, \
         patch('file_service.service.files_collection') as mock_files:

        mock_files.count_documents.return_value = 0
        mock_s3.create_multipart_upload.return_value = {
            'UploadId': 'test_upload_id_123'
        }
        mock_sessions.insert_one.return_value = Mock(inserted_id=ObjectId())

        response = file_service.InitiateMultipartUpload(request, mock_context)

        assert response.upload_id == 'test_upload_id_123'
        assert response.chunk_size == 10 * 1024 * 1024
        assert response.total_parts == 10

        mock_s3.create_multipart_upload.assert_called_once()
        mock_sessions.insert_one.assert_called_once()

def test_upload_part(file_service, mock_context):
    request = file_pb2.UploadPartRequest(
        upload_id="test_upload_id_123",
        part_number=1,
        chunk=b"x" * (10 * 1024 * 1024)
    )

    with patch('file_service.service.s3_client') as mock_s3, \
         patch('file_service.service.upload_sessions_collection') as mock_sessions:

        mock_s3.upload_part.return_value = {
            'ETag': '"abc123etag"'
        }
        mock_sessions.find_one.return_value = {
            "upload_id": "test_upload_id_123",
            "user_id": "test_user_123",
            "s3_key": "test_user_123/file_id/file.mp4",
            "parts": []
        }
        mock_sessions.update_one.return_value = Mock()

        response = file_service.UploadPart(request, mock_context)

        assert response.etag == '"abc123etag"'
        assert response.part_number == 1

        mock_s3.upload_part.assert_called_once()
        mock_sessions.update_one.assert_called_once()

def test_complete_multipart_upload(file_service, mock_context):
    file_id = str(ObjectId())
    parts = [
        file_pb2.PartInfo(part_number=1, etag='"etag1"'),
        file_pb2.PartInfo(part_number=2, etag='"etag2"'),
    ]
    request = file_pb2.CompleteMultipartUploadRequest(
        upload_id="test_upload_id_123",
        parts=parts
    )

    with patch('file_service.service.s3_client') as mock_s3, \
         patch('file_service.service.upload_sessions_collection') as mock_sessions, \
         patch('file_service.service.files_collection') as mock_files:

        mock_sessions.find_one.return_value = {
            "upload_id": "test_upload_id_123",
            "user_id": "test_user_123",
            "file_id": file_id,
            "filename": "large_file.mp4",
            "content_type": "video/mp4",
            "total_size": 20 * 1024 * 1024,
            "s3_key": "test_user_123/file_id_123/large_file.mp4"
        }
        mock_s3.complete_multipart_upload.return_value = {}
        mock_files.insert_one.return_value = Mock(inserted_id=ObjectId(file_id))
        mock_sessions.delete_one.return_value = Mock()

        response = file_service.CompleteMultipartUpload(request, mock_context)

        assert response.file.id == file_id
        assert response.file.filename == "large_file.mp4"
        assert response.file.user_id == "test_user_123"

        mock_s3.complete_multipart_upload.assert_called_once()
        mock_files.insert_one.assert_called_once()
        mock_sessions.delete_one.assert_called_once()

def test_abort_multipart_upload(file_service, mock_context):
    request = file_pb2.AbortMultipartUploadRequest(
        upload_id="test_upload_id_123"
    )

    with patch('file_service.service.s3_client') as mock_s3, \
         patch('file_service.service.upload_sessions_collection') as mock_sessions:

        mock_sessions.find_one.return_value = {
            "upload_id": "test_upload_id_123",
            "user_id": "test_user_123",
            "s3_key": "test_user_123/file_id/file.mp4"
        }
        mock_s3.abort_multipart_upload.return_value = {}
        mock_sessions.delete_one.return_value = Mock(deleted_count=1)

        response = file_service.AbortMultipartUpload(request, mock_context)

        assert response.success is True

        mock_s3.abort_multipart_upload.assert_called_once()
        mock_sessions.delete_one.assert_called_once()
