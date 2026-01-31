from pymongo import MongoClient
from internal.config.config import MONGO_URI

client = MongoClient(MONGO_URI)
db = client["file_service"]
files_collection = db["files"]