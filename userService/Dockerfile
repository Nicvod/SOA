FROM golang:1.23

WORKDIR /app

COPY service /app/service
COPY .keys /app/.keys
COPY user_proto /app/user_proto

WORKDIR /app/service

RUN go build -o main

EXPOSE 50051

CMD ["./main", "-private_key=/app/.keys/signature.pem", "-public_key=/app/.keys/signature.pub", "-db_name_env=POSTGRES_DB", "-db_user_env=POSTGRES_USER", "-db_password_env=POSTGRES_PASSWORD", "-db_port=5432", "-service_port=50051"]