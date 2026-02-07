from unittest.mock import Mock, patch


def test_serve_wires_grpc_and_telemetry_without_running_real_server():
    server = Mock()

    with patch("file_service.__main__.init_telemetry") as init_telemetry, patch(
        "file_service.__main__.grpc.server", return_value=server
    ) as grpc_server, patch(
        "file_service.__main__.futures.ThreadPoolExecutor"
    ) as executor_cls, patch(
        "file_service.__main__.file_pb2_grpc.add_FileServiceServicer_to_server"
    ) as add_servicer, patch(
        "file_service.__main__.register_health"
    ) as register_health, patch(
        "file_service.__main__.SERVICE_PORT", 55555
    ), patch(
        "file_service.__main__.SERVICE_NAME", "file-service"
    ), patch(
        "file_service.__main__.ENVIRONMENT", "test"
    ), patch(
        "file_service.__main__.OTLP_ENDPOINT", "http://otel"
    ), patch(
        "file_service.__main__.AXIOM_TOKEN", "token"
    ), patch(
        "file_service.__main__.DATASET", "dataset"
    ):
        from file_service.__main__ import serve

        # Prevent blocking forever.
        server.wait_for_termination.side_effect = SystemExit(0)

        try:
            serve()
        except SystemExit:
            pass

        init_telemetry.assert_called_once_with(
            service_name="file-service",
            environment="test",
            endpoint="http://otel",
            token="token",
            dataset="dataset",
            metrics_dataset="metrics",
        )

        executor_cls.assert_called_once_with(max_workers=10)
        grpc_server.assert_called_once()

        add_servicer.assert_called_once()
        register_health.assert_called_once_with(server)

        server.add_insecure_port.assert_called_once_with("[::]:55555")
        server.start.assert_called_once_with()
        server.wait_for_termination.assert_called_once_with()
