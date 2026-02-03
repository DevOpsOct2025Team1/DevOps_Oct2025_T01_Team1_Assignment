import os
from pymongo import MongoClient
from file_service.config import MONGO_URI

client = MongoClient(MONGO_URI)
DB_NAME = os.getenv("MONGODB_DATABASE", "file_service")
db = client[DB_NAME]
files_collection = db["files"]