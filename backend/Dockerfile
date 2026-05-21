FROM golang:1.25.3 AS builder

WORKDIR /app

COPY . .

RUN cd cmd && CGO_ENABLED=0 GOOS=linux go build -o app main.go

FROM alpine AS app

WORKDIR /app

COPY --from=builder /app/cmd/app ./
COPY --from=builder /app/env ./env
COPY --from=builder /app/entrypoint.sh ./entrypoint.sh
RUN chmod +x ./entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]
