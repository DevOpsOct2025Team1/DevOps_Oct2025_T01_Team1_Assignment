import importlib
from unittest.mock import MagicMock, Mock, patch


def test_store_uses_env_database_and_mongo_uri(monkeypatch):
    monkeypatch.setenv("MONGO_URI", "mongodb://example:27017")
    monkeypatch.setenv("MONGODB_DATABASE", "test_db")

    mock_client = MagicMock()
    mock_db = MagicMock()
    mock_collection = Mock()

    mock_client.__getitem__.return_value = mock_db
    mock_db.__getitem__.return_value = mock_collection

    # Patch pymongo.MongoClient BEFORE reloading the module so module-level
    # initialization doesn't attempt a real connection.
    with patch("pymongo.MongoClient", return_value=mock_client) as mongo_cls:
        import file_service.config as config
        import file_service.store as store

        importlib.reload(config)
        reloaded = importlib.reload(store)

        mongo_cls.assert_called_once_with("mongodb://example:27017")
        assert reloaded.DB_NAME == "test_db"
        assert reloaded.files_collection is mock_collection
