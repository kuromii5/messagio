# Messagio gRPC microservice with Kafka

Messagio is gRPC microservice on Go with gRPC gateway handling HTTP requests which manages message sending and getting using Apache Kafka. It also uses PostgreSQL as a database.

## Installation

First you need to configure .env file. Here is an example of configuration:

```env
# ENVIRONMENT
ENV=local

# LOG LEVEL
LOG=info

# SERVER
GRPC_PORT=44044
HTTP_PORT=8080

# POSTGRES SETTINGS
POSTGRES_USER=postgres
POSTGRES_PASSWORD=admin
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_DBNAME=messagio
POSTGRES_SSLMODE=disable

# KAFKA
KAFKA_BROKERS=kafka:9092
KAFKA_PORT=9092
KAFKA_TOPIC=main # set only one topic :(
KAFKA_BROKER_ID=1
KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181
KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://kafka:9092
KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1

# ZOOKEEPER
ZOOKEEPER_CLIENT_PORT=2181
ZOOKEEPER_TICK_TIME=2000
```

After creating configuration you can run next command:

```bash
docker-compose up --build
```

All necessary stuff will be automatically generated in docker-containers and application will run.

## Usage

I think it's too expensive to generate a Swagger documentation for it, so...

Using Postman and similar tools you can send next requests to this little API, for example:

```bash
(GET): http://localhost:8080/v1/stats # allows you to see amount of all processed messages
(GET): http://localhost:8080/v1/get_messages # consumes all new messages from Kafka
(POST): http://localhost:8080/v1/send_message # creates message to save in DB and send it to Kafka
{
    "message":"Hello World!" # the only field you need to specify
}
```
