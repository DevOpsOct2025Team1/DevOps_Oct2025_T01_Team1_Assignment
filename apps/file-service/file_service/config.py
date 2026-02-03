import os
from dotenv import load_dotenv

load_dotenv()

MONGO_URI = os.getenv("MONGO_URI", "mongodb://localhost:27017")
SERVICE_PORT = int(os.getenv("FILE_SERVICE_PORT", 50054))
SERVICE_NAME = "file-service"
ENVIRONMENT = os.getenv("ENVIRONMENT", "development")
OTLP_ENDPOINT = os.getenv("OTLP_ENDPOINT", "https://us-east-1.aws.edge.axiom.co")
AXIOM_TOKEN = os.getenv("AXIOM_API_TOKEN", "")
DATASET = os.getenv("AXIOM_DATASET", "")