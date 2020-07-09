# MonoBank -  telegram chat bot

Loads transactions from Monobank and prepares data for the [Money Pro](https://money.pro/mac/) application. You will have the opportunity to receive data with a large number of filters.
## Precondition

* Installed [Docker](https://docs.docker.com/engine/install/).

## Installation

```bash
docker build -t "tag_name" .
```

## Environment Variables for the Docker Container
* TOKEN - token for connecting to Telegram
* TIMEOUT - Telegram offset update
* REDIS_URL - url for connecting to Redis

## Test
* Run tests.
```bash
make test
```

## Linters
* Run linters.
```bash
make lint
```