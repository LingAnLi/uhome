# PostHouseImg Service

This is the PostHouseImg service

Generated with

```
micro new uhone/PostHouseImg --namespace=go.micro --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.srv.PostHouseImg
- Type: srv
- Alias: PostHouseImg

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
./PostHouseImg-srv
```

Build a docker image
```
make docker
```