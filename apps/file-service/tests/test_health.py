from unittest.mock import Mock, patch

from file_service.health import register_health


def test_register_health_adds_servicer_to_server():
    server = Mock()

    with patch("grpc_health.v1.health.HealthServicer") as health_servicer_cls, patch(
        "grpc_health.v1.health_pb2_grpc.add_HealthServicer_to_server"
    ) as add_to_server:
        health_servicer = Mock()
        health_servicer_cls.return_value = health_servicer

        register_health(server)

        health_servicer_cls.assert_called_once_with()
        add_to_server.assert_called_once_with(health_servicer, server)
