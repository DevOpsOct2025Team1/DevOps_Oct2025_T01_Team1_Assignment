FROM python:3.11-slim AS builder

WORKDIR /build

COPY requirements.txt .

RUN pip install --upgrade pip \
    && pip install --prefix=/install -r requirements.txt

COPY . .

FROM python:3.11-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /apps

COPY --from=builder /install /usr/local
COPY --from=builder /build .

EXPOSE 50054

ENV PYTHONPATH=/apps

# Run the server
CMD ["python", "-u", "-m", "cmd.server.main"]