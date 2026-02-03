FROM ghcr.io/astral-sh/uv:python3.11-bookworm-slim AS builder

WORKDIR /build

COPY pyproject.toml .
COPY uv.lock* .

RUN uv sync --frozen --no-dev --no-install-project

COPY file_service ./file_service
COPY user ./user
COPY auth ./auth
COPY file ./file

RUN uv sync --frozen --no-dev

FROM python:3.11-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /build/.venv /app/.venv
COPY --from=builder /build/file_service ./file_service
COPY --from=builder /build/user ./user
COPY --from=builder /build/auth ./auth
COPY --from=builder /build/file ./file

EXPOSE 50054

ENV PATH="/app/.venv/bin:$PATH"
ENV PYTHONPATH=/app

CMD ["python", "-u", "-m", "file_service"]