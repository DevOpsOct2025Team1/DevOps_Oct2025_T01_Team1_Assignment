import os
from concurrent import futures
import grpc

from file_service.service import FileService
from file_service.auth_client import AuthClient
from file.v1 import file_pb2_grpc
from file_service.health import register_health
from file_service.config import FILE_SERVICE_PORT, SERVICE_NAME, ENVIRONMENT, OTLP_ENDPOINT, AXIOM_TOKEN, DATASET
SERVICE_PORT = FILE_SERVICE_PORT
from file_service.telemetry import init_telemetry

def serve():
    init_telemetry(
        service_name=SERVICE_NAME,
        environment=ENVIRONMENT,
        endpoint=OTLP_ENDPOINT,
        token=AXIOM_TOKEN,
        dataset=DATASET,
        metrics_dataset="metrics",
    )

    auth_client = AuthClient()

    max_msg_size = 20 * 1024 * 1024
    server = grpc.server(
        futures.ThreadPoolExecutor(max_workers=10),
        options=[
            ("grpc.max_receive_message_length", max_msg_size),
            ("grpc.max_send_message_length", max_msg_size),
        ],
    )
    file_pb2_grpc.add_FileServiceServicer_to_server(FileService(auth_client), server)
    register_health(server)

    server.add_insecure_port(f"[::]:{FILE_SERVICE_PORT}")
    print(f"File service gRPC running on port {FILE_SERVICE_PORT}")
    server.start()
    try:
        server.wait_for_termination()
    finally:
        auth_client.close()


if __name__ == "__main__":
    serve()