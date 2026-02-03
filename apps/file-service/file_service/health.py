from grpc_health.v1 import health_pb2, health_pb2_grpc


def register_health(server):
    health_servicer = health_pb2_grpc.HealthServicer()
    health_pb2_grpc.add_HealthServicer_to_server(health_servicer, server)