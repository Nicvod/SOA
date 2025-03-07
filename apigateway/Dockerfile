FROM golang:1.23

WORKDIR /app

COPY ./apigateway/service /app/service
COPY ./userService/user_proto /userService/user_proto

WORKDIR /app/service

RUN go build -o main

EXPOSE 50051

CMD ["./main", "--rest_endpoint=80", "--user_service_endpoint=50051"]