# Signature Service

Demo application to create signature devices. With the created signature devices the signature service can also create signatures for the provided data.

## Prerequisites & Tooling

- Golang (v1.21+)

## Usage

### Start the app
```sh
$ make run
```

### Run the tests
```sh
$ make test
```

### Create signature device

```sh
curl --location 'localhost:8000/api/v1/device' \
--header 'Content-Type: text/plain' \
--data '{"id": "bc890106-641e-41fc-aed8-b77cca0b42b9", "algorithm": "RSA", "label": "NEW"}'
```

### Sign data

```sh
curl --location 'localhost:8000/api/v1/signature' \
--header 'Content-Type: text/plain' \
--data '{"device_id": "bc890106-641e-41fc-aed8-b77cca0b42b9", "data_to_be_signed": "datatosign"}'
```

## Architecture

| Package | Description |
| ------ | ------ |
| crypt | Utilities to deal with cryptography |
| device | Signature device domain |
| signature | Signature domain |
| health | Contains health check handler |
| api | Contains http api related logics |