# MVP-proxy service

This is a Go-based proxy service designed to handle both HTTP and HTTPS connections. It includes user authentication and bandwidth limitation features, with hardcoded credentials for an MVP version. The service also exposes several statistics endpoints.

## Table of contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
- [Endpoints](#endpoints)
- [Authentication](#authentication)
- [Configuration](#configuration)
- [Running tests](#running-tests)
- [Makefile commands](#makefile-commands)
- [Deployment](#deployment)
- [TODO](#todo)

## Features

- **Proxy for HTTP and HTTPS**: Handles secure and non-secure traffic.
- **Basic authentication**: Hardcoded credentials for user authentication.
- **Bandwidth limitation**: 1GB bandwidth usage limit per user.
- **Statistics endpoints**: Check the health, bandwidth usage, and history of accessed URLs.
- **Request duration limit**: 2-hour request limit, configurable in the config file.

## Requirements

- **Docker**: Version 27.0+ (for containerized usage).
- **Go**: Version 1.22+ (if running locally without Docker).
- **Makefile**: Simplifies common tasks like running the app and debugging.

## Installation

### Running with Docker

1. Ensure you have Docker installed.
2. Clone the repository:

   ```bash
   git clone https://github.com/krizanauskas/mvp-proxy.git
   cd proxy-service
   ```

3. Build and run the Docker container:

   ```bash
   docker-compose up
   ```

### Running without Docker

1. Install Go (version 1.22+).
2. Clone the repository:

   ```bash
   git clone https://github.com/krizanauskas/mvp-proxy.git
   cd proxy-service
   ```

3. Set up environment variables by creating a `.env` file:  
    Copy the contents of `.env.example` to `.env`.

    More on this in the configuration section.

4. Run the service using the Makefile:

   ```bash
   make proxy-run
   ```

## Usage

Once the proxy is running, it listens for traffic on `localhost:8080` (or a custom port based on the config).

Example request using curl:

```bash
curl -v -x http://127.0.0.1:8080 --proxy-user user:pass https://www.google.com
```

### Authentication

The service uses basic authentication. For the MVP version, the following hardcoded credentials are used:

- **Username**: `user`
- **Password**: `pass`

## Endpoints

In addition to the proxy functionality, the service provides the following endpoints for statistics:

1. **Health check**: `GET http://localhost:3333/health`  
   - Returns the health status of the proxy server.
   
2. **Usage limits**: `GET http://localhost:3333/usage-limits`  
   - Shows the bandwidth limit and the current usage for the authenticated user (1GB limit).

3. **History**: `GET http://localhost:3333/history`  
   - Displays a list of previously accessed URLs through the proxy.

## Configuration

The application configuration can be managed through environment variables and YAML config files. By default, the application reads from the `config/config-dev.yaml` file. You can create a custom configuration by setting the `APP_ENV` variable and adding a corresponding YAML file in the `config` directory.

### Example for development configuration (`config/config-dev.yaml`):

```yaml
proxy_server:
  port: ":8080"
  max_request_duration_sec: 7200
  allowed_data_mb: 1000

status_server:
  port: ":3333"
```

- **Example for local configuration**:

1. Set `APP_ENV=local` in your `.env` file.
   
2. Create a `config-local.yaml` file in the `config` directory:

   ```yaml
   proxy_server:
     port: ":8081"
     max_request_duration_sec: 7200
     allowed_data_mb: 1000

   status_server:
     port: ":3334"
   ```

This will change the proxy port to `8081` and the stats port to `3334`.

## Running tests

You can run tests for the application locally using the following command:

```bash
go test ./...
```

## Makefile commands

The `Makefile` simplifies common tasks. Available commands include:

- **proxy-run**: Runs the proxy service.
  
  ```bash
  make proxy-run
  ```

- **proxy-debug**: Starts a debugging session using `dlv` on port `2345`.
  
  ```bash
  make proxy-debug
  ```

- **proxy-race**: Runs the proxy with the race detector enabled.
  
  ```bash
  make proxy-race
  ```

## Deployment

The project is configured to deploy via **GitHub Actions** using a **self-hosted runner**. The pipeline does the following:

1. Checks out the code.
2. Builds a Docker image using the provided `Dockerfile`.
3. Pushes the image to **GitHub Container Registry**.
4. Updates the Docker Swarm service with the new image.

## TODO

- **Track uploaded HTTP request traffic**: Implement functionality to track the traffic uploaded through the proxy.
- **Fix random HTTP errors**: Investigate and fix random HTTP errors currently being thrown by the proxy.
- **Improve error handling and logging**: Enhance error handling mechanisms and provide more detailed logging to improve debugging and operational transparency.