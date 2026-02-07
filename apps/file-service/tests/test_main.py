from unittest.mock import Mock, patch


def test_serve_wires_grpc_and_telemetry_without_running_real_server():
    server = Mock()

    import file_service.__main__ as main

    with patch.object(main, "init_telemetry") as init_telemetry, patch.object(
        main.grpc, "server", return_value=server
    ) as grpc_server, patch.object(
        main.futures, "ThreadPoolExecutor"
    ) as executor_cls, patch.object(
        main.file_pb2_grpc, "add_FileServiceServicer_to_server"
    ) as add_servicer, patch.object(
        main, "register_health"
    ) as register_health, patch.object(
        main, "SERVICE_PORT", 55555
    ), patch.object(
        main, "SERVICE_NAME", "file-service"
    ), patch.object(
        main, "ENVIRONMENT", "test"
    ), patch.object(
        main, "OTLP_ENDPOINT", "http://otel"
    ), patch.object(
        main, "AXIOM_TOKEN", "token"
    ), patch.object(
        main, "DATASET", "dataset"
    ):

        # Prevent blocking forever.
        server.wait_for_termination.side_effect = SystemExit(0)

        try:
            main.serve()
        except SystemExit:
            # Swallow the expected SystemExit from wait_for_termination so the
            # test can continue with assertions.
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
