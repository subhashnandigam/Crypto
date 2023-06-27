
# Project CRYPTO

This repository contains a Go project for managing cryptocurrency data. 


## Project Structure

The project's structure is organized as follows:

```
.
├── app
│   ├── routes.go
│   └── syncmap.go
├── go.mod
├── go.sum
└── main.go

```

## Running the Application

To run the application, navigate to the project's root directory:

```bash
cd E:\go-workspace\src\github.com\crypto
```

Then, execute the following command:

```bash
go run cmd/main.go
```

The application will start, and you can access the following endpoints:

- GET http://localhost:8000/currency/all
- POST http://localhost:8000/internal/add/support/TLMUSDT

Feel free to modify the host and port settings in the `main.go` file according to your preferences.

## Dependencies

The project manages its dependencies using Go modules. The necessary dependencies are specified in the `go.mod` file. When you run the application for the first time, Go will automatically download the required dependencies.

If you need to add or update dependencies, you can use the `go get` command followed by the package name. For example:


