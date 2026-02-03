from concurrent import futures
import grpc
from internal.service.service import FileService
from file.v1 import file_pb2_grpc
from internal.health.health import register_health
from internal.config.config import SERVICE_PORT
from telemetry import init_telemetry

SERVICE_NAME = "file-service"
ENVIRONMENT = "development"
OTLP_ENDPOINT = "https://us-east-1.aws.edge.axiom.co"
AXIOM_TOKEN = "YOUR_AXIOM_TOKEN"
DATASET = "your_dataset"

def serve():
    # Initialize telemetry first
    init_telemetry(
        service_name=SERVICE_NAME,
        environment=ENVIRONMENT,
        endpoint=OTLP_ENDPOINT,
        token=AXIOM_TOKEN,
        dataset=DATASET,
        metrics_dataset="metrics",
    )

    # Start gRPC server
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    file_pb2_grpc.add_FileServiceServicer_to_server(FileService(), server)
    register_health(server)

    server.add_insecure_port(f"[::]:{SERVICE_PORT}")
    print(f"File service running on port {SERVICE_PORT}")
    server.start()
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
