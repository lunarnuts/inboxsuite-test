FROM golang:alpine3.20 as builder
WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o main ./cmd/main \
    && CGO_ENABLED=0 GOOS=linux go build -a -o migrator ./cmd/migrator \
    && mkdir final_dir \
    && mv main migrator final_dir/


FROM alpine:latest
WORKDIR /app/

COPY --from=builder /app/final_dir .

CMD /app/migrator -config /tmp/config.yaml ; /app/main -config /tmp/config.yaml
