# angkutgan
Based on Scaffolding Samarinda Project. 
- GoLang as Backend
- jQuery as Frontend
- MySQL as Database
- Linux as Compatible Development Environment

## Getting Started
- run command `make config`
- fill necessary data to `config.yaml` (jwt key, db configuration)

## Run Service
- make sure `config.yaml` is ready
- run command `go mod vendor`
- run command `go run .`

## Deployment
- make sure `config.yaml` is ready
- run command `make build`
- run one of binary from `executable` directory based on environment, run `make check` to make sure the current environment.
