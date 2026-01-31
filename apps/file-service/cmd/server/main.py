from concurrent import futures
import grpc
from internal.service.service import FileService
from proto.file.v1 import file_pb2_grpc
from internal.health.health import register_health
from internal.config.config import SERVICE_PORT
from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor, ConsoleSpanExporter


def setup_tracing():
    trace.set_tracer_provider(TracerProvider())
    tracer = trace.get_tracer(__name__)
    span_processor = BatchSpanProcessor(ConsoleSpanExporter())
    trace.get_tracer_provider().add_span_processor(span_processor)
    return tracer

def serve():
    tracer = setup_tracing()

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

    file_pb2_grpc.add_FileServiceServicer_to_server(FileService(), server)
    register_health(server)

    server.add_insecure_port(f"[::]:{SERVICE_PORT}")
    print(f"File service running on port {SERVICE_PORT}")
    server.start()
    server.wait_for_termination()

if __name__ == "__main__":
    serve()