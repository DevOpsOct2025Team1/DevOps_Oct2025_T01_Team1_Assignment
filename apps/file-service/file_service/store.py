import os
import boto3
from pymongo import MongoClient
from file_service.config import MONGO_URI, S3_ENDPOINT, S3_ACCESS_KEY, S3_SECRET_KEY, S3_BUCKET_NAME

client = MongoClient(MONGO_URI)
DB_NAME = os.getenv("MONGODB_DATABASE", "file_service")
db = client[DB_NAME]
files_collection = db["files"]

# Initialize S3 client
s3_config = {
    'service_name': 's3',
    'region_name': 'us-east-1'
}

if S3_ENDPOINT:
    s3_config['endpoint_url'] = S3_ENDPOINT

if S3_ACCESS_KEY:
    s3_config['aws_access_key_id'] = S3_ACCESS_KEY

if S3_SECRET_KEY:
    s3_config['aws_secret_access_key'] = S3_SECRET_KEY

s3_client = boto3.client(**s3_config)

def generate_s3_key(user_id: str, file_id: str, filename: str) -> str:
    """Generate unique S3 key for file storage."""
    return f"{user_id}/{file_id}/{filename}"