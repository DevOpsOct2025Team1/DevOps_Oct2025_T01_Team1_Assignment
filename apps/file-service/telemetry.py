# telemetry.py
from opentelemetry import trace, metrics
from opentelemetry.sdk.resources import Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.sdk.metrics import MeterProvider
from opentelemetry.sdk.metrics.export import PeriodicExportingMetricReader
from opentelemetry.exporter.otlp.proto.http.trace_exporter import OTLPSpanExporter
from opentelemetry.exporter.otlp.proto.http.metric_exporter import OTLPMetricExporter
from opentelemetry.propagate import set_global_textmap
from opentelemetry.propagators.composite import CompositePropagator
from opentelemetry.trace.propagation.tracecontext import TraceContextTextMapPropagator
from opentelemetry.instrumentation.requests import RequestsInstrumentor


def init_telemetry(
    service_name: str,
    environment: str = "development",
    endpoint: str = "https://us-east-1.aws.edge.axiom.co",
    token: str = None,
    dataset: str = None,
    metrics_dataset: str = "metrics",
):
    if not service_name or not token or not dataset:
        raise ValueError("service_name, token, and dataset are required")

    resource = Resource.create({"service.name": service_name, "environment": environment})


    trace_headers = {
        "Authorization": f"Bearer {token}",
        "X-AXIOM-DATASET": dataset,
    }
    trace_exporter = OTLPSpanExporter(endpoint=endpoint, headers=trace_headers)
    tracer_provider = TracerProvider(resource=resource)
    tracer_provider.add_span_processor(BatchSpanProcessor(trace_exporter))
    trace.set_tracer_provider(tracer_provider)


    metrics_headers = {
        "Authorization": f"Bearer {token}",
        "X-AXIOM-DATASET": metrics_dataset,
    }
    metric_exporter = OTLPMetricExporter(endpoint=endpoint, headers=metrics_headers)
    meter_provider = MeterProvider(
        resource=resource,
        metric_readers=[PeriodicExportingMetricReader(metric_exporter)]
    )
    metrics.set_meter_provider(meter_provider)


    propagator = TraceContextTextMapPropagator()
    set_global_textmap(propagator)

    RequestsInstrumentor().instrument()

    return trace.get_tracer(__name__), meter_provider
