import pytest

from file_service.telemetry import init_telemetry


def test_init_telemetry_requires_service_name_token_dataset():
    with pytest.raises(ValueError):
        init_telemetry(service_name="")

    with pytest.raises(ValueError):
        init_telemetry(service_name="file-service", token="", dataset="ds")

    with pytest.raises(ValueError):
        init_telemetry(service_name="file-service", token="t", dataset="")
