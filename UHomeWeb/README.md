# UHomeWeb Service

This is the UHomeWeb service

Generated with

```
micro new uhone/UHomeWeb --namespace=go.micro --type=web
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.web.UHomeWeb
- Type: web
- Alias: UHomeWeb

## Dependencies

Micro services depend on service discovery. The default is consul.

```
# install consul
brew install consul

# run consul
consul agent -dev
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./UHomeWeb-web
```

Build a docker image
```
make docker
```