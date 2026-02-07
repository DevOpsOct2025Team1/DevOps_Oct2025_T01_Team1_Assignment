import os
from dotenv import load_dotenv

load_dotenv()

MONGO_URI = os.getenv("MONGO_URI", "mongodb://localhost:27017")
FILE_SERVICE_PORT = os.getenv("FILE_SERVICE_PORT", "50054")
AUTH_SERVICE_ADDR = os.getenv("AUTH_SERVICE_ADDR", "localhost:8081")
HTTP_PORT = os.getenv("HTTP_PORT", "3001")

# S3 Configurationq
S3_ENDPOINT = os.getenv("S3_ENDPOINT", "https://clw-s3.ngeeann.zip")
S3_ACCESS_KEY = os.getenv("S3_ACCESS_KEY", "")
S3_SECRET_KEY = os.getenv("S3_SECRET_KEY", "")
S3_BUCKET_NAME = os.getenv("S3_BUCKET_NAME", "file-storage")
SERVICE_NAME = "file-service"
ENVIRONMENT = os.getenv("ENVIRONMENT", "development")
OTLP_ENDPOINT = os.getenv("OTLP_ENDPOINT", "https://us-east-1.aws.edge.axiom.co")
AXIOM_TOKEN = os.getenv("AXIOM_API_TOKEN", "")
DATASET = os.getenv("AXIOM_DATASET", "")