from pymongo import MongoClient
from internal.config.config import MONGODB_URI, MONGODB_DATABASE

client = MongoClient(MONGODB_URI)
db = client[MONGODB_DATABASE]
files_collection = db["files"]