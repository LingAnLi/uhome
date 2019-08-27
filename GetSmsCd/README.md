# GetSmsCd Service

This is the GetSmsCd service

Generated with

```
micro new uhone/GetSmsCd --namespace=go.micro --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.srv.GetSmsCd
- Type: srv
- Alias: GetSmsCd

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
./GetSmsCd-srv
```

Build a docker image
```
make docker
```