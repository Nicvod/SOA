FROM golang:1.23

WORKDIR /app

COPY . /app
COPY ./userService/.keys /app/.keys

WORKDIR /app

RUN go build -o main ./postService/cmd/...

EXPOSE 50051

CMD ["./main", "-public_key=/app/.keys/signature.pub", "-db_name_env=POSTGRES_DB", "-db_user_env=POSTGRES_USER", "-db_password_env=POSTGRES_PASSWORD", "-db_port=5432", "-service_port=50051"]