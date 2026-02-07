import grpc
from auth.v1 import auth_pb2, auth_pb2_grpc
from file_service.config import AUTH_SERVICE_ADDR

class AuthClient:
    def __init__(self):
        self.channel = grpc.insecure_channel(AUTH_SERVICE_ADDR)
        self.stub = auth_pb2_grpc.AuthServiceStub(self.channel)

    def validate_token(self, token: str) -> auth_pb2.ValidateTokenResponse:
        try:
            return self.stub.ValidateToken(auth_pb2.ValidateTokenRequest(token=token))
        except grpc.RpcError as e:
            print(f"Auth verification failed: {e}")
            return None

    def close(self):
        self.channel.close()
