import os
from dotenv import load_dotenv

load_dotenv()

MONGO_URI = os.getenv("MONGO_URI", "mongodb://localhost:27017")
SERVICE_PORT = int(os.getenv("FILE_SERVICE_PORT", 50054))