import os
from concurrent import futures
import grpc

from file_service.service import FileService
from file.v1 import file_pb2_grpc
from file_service.health import register_health
from file_service.config import FILE_SERVICE_PORT, SERVICE_NAME, ENVIRONMENT, OTLP_ENDPOINT, AXIOM_TOKEN, DATASET
SERVICE_PORT = FILE_SERVICE_PORT
from file_service.telemetry import init_telemetry

import threading
from file_service.http_server import run_http_server

def serve():
    init_telemetry(
        service_name=SERVICE_NAME,
        environment=ENVIRONMENT,
        endpoint=OTLP_ENDPOINT,
        token=AXIOM_TOKEN,
        dataset=DATASET,
        metrics_dataset="metrics",
    )

    # Start HTTP server in a separate thread
    if os.getenv("FILE_SERVICE_ENABLE_HTTP", "1") == "1":
        http_thread = threading.Thread(target=run_http_server, daemon=True)
        http_thread.start()

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    file_pb2_grpc.add_FileServiceServicer_to_server(FileService(), server)
    register_health(server)

    server.add_insecure_port(f"[::]:{FILE_SERVICE_PORT}")
    print(f"File service gRPC running on port {FILE_SERVICE_PORT}")
    server.start()
    server.wait_for_termination()


if __name__ == "__main__":
    serve()