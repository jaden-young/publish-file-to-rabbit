# Publish File to Rabbit

Reads a given file with a single array of JSON objects and publishes each object
as a single message to a RabbitMQ queue.

## Usage

### With Go

Explicitly passing default options:

```sh
go get -u github.com/jaden-young/publish-file-to-rabbit
publish-file-to-rabbit --host localhost --port 5672 --user guest \
  --password guest --queue eiffel --file events.json
```

### With Docker

Explicitly passing default options (other than the volume):

```sh
docker run -it --rm \
  -v data:data \
  --env AMQP_HOST=localhost \
  --env AMQP_PORT=5672 \
  --env AMQP_USER=guest \
  --env AMQP_PASSWORD=guest \
  --env QUEUE_NAME=eiffel \
  --env EVENTS_FILE=events.json \
  jadyoung/publish-file-to-rabbit
```

The docker image contains a default file of json objects to publish,
and also will wait for the RabbitMQ host to be available before publishing.