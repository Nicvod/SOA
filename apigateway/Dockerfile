FROM golang:1.23

WORKDIR /app

COPY . /app
COPY ./userService/.keys /app/.keys

WORKDIR /app

RUN go build -o main ./apigateway/service/...

EXPOSE 50051

CMD ["./main", "--rest_endpoint=80", "--user_service_endpoint=user_app:50051", "--post_service_endpoint=post_app:50051", "--public_key=/app/.keys/signature.pub"]