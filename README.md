# groxy

A simple proxy to test stuff

## Usage

```sh
docker pull registry.fabiomora.dev/groxy:latest
```

### docker compose

```yml
services:
    example-service:
        image: registry.fabiomora.dev/groxy:latest
        command: --config /etc/groxy/groxy.yaml
        volumes:
            - ./config.yml:/etc/groxy/groxy.yaml
```

> [Example configuration](example.yml)

## Demo

![Demonstration](assets/example.gif)
