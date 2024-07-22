FROM golang:1.22

# Install protobuf compiler
RUN apt-get update && DEBIAN_FRONTEND=nointeractive apt-get install --no-install-recommends --assume-yes protobuf-compiler

WORKDIR /app
COPY . .
RUN go mod download && go mod verify

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2

# Run the Makefile build script to install dependencies,
# get proto files, generate code, and run migrations
RUN make build

# Run migrations and start the application
CMD ["sh", "-c", "make migrate-up && make run"]