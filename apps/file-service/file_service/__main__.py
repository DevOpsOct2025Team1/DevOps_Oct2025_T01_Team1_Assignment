import os
from concurrent import futures
import grpc

from file_service.service import FileService
from file.v1 import file_pb2_grpc
from file_service.health import register_health
from file_service.config import SERVICE_PORT, SERVICE_NAME, ENVIRONMENT, OTLP_ENDPOINT, AXIOM_TOKEN, DATASET
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

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    file_pb2_grpc.add_FileServiceServicer_to_server(FileService(), server)
    register_health(server)

    server.add_insecure_port(f"[::]:{SERVICE_PORT}")
    print(f"File service running on port {SERVICE_PORT}")
    server.start()
    server.wait_for_termination()


if __name__ == "__main__":
    serve()