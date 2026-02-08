FROM cgr.dev/chainguard/python:latest-dev AS builder

ENV LANG=C.UTF-8
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1

WORKDIR /build

RUN pip install --no-cache-dir uv

COPY apps/file-service/pyproject.toml .
COPY apps/file-service/uv.lock* .

RUN uv export --frozen --no-dev --no-emit-project -o requirements.txt && \
    uv pip install --no-cache-dir --target /build/deps -r requirements.txt

COPY apps/file-service/file_service ./file_service
COPY apps/file-service/user ./user
COPY apps/file-service/auth ./auth
COPY apps/file-service/file ./file

FROM cgr.dev/chainguard/python:latest

ENV PYTHONUNBUFFERED=1
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONPATH=/app/deps:/app

WORKDIR /app

COPY --from=builder /build/deps /app/deps
COPY --from=builder /build/file_service ./file_service
COPY --from=builder /build/user ./user
COPY --from=builder /build/auth ./auth
COPY --from=builder /build/file ./file

EXPOSE 50054

CMD ["-u", "-m", "file_service"]
