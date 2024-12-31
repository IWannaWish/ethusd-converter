# gRPC API Template

This repository provides a template for implementing a gRPC API in Go, which can also be accessed via REST. It demonstrates a well-structured project layout that adheres to best practices in Go development.

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Building and Running](#building-and-running)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## Features

- **gRPC API**: Fast and efficient remote procedure calls.
- **RESTful Access**: Access the same functionality through RESTful endpoints.
- **Best Practices**: Organized project structure for maintainability and scalability.
- **Integrated Testing**: Unit and integration tests to ensure code quality.
- **Protobuf Definitions**: Clear and structured API definitions using Protocol Buffers.
- s fh
## Getting Started

### Prerequisites

- Go 1.18 or later
- Protobuf compiler (`protoc`)
  ```bash
  brew install protobuf # for mac
  ```
    ```bash
  pip install protobuf # for linux
  ```
- gRPC and related Go libraries
  ```bash
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
  ```
- Pre-commit
  ```bash
  brew install pre-commit # for mac
  ```
  ```bash
  pip install pre-commit # for linux
  ```

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/TimRutte/api.git
   cd api
   ```

2. Install dependencies:

   ```bash
   make tidy
   ```

3. Install necessary Go tools:

   ```bash
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
   ```

## Building and Running

To build and run the application, use the following commands:

```bash
make build-linux
```
```bash
make build-darwin
```
```bash
make build-windows
```


```bash
make run
```

### Docker

To run the application in a Docker container, build the image:

```bash
make docker-build
```

Then run the container:

```bash
make docker-run
```

## Testing

Unit tests and integration tests are included in the project. To run the tests with vulnerability check, use:

```bash
make test
```

To run only the unit tests:

```bash
make test-only
```

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
