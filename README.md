# RoastGitHub API

RoastGitHub API is a service that interacts with GitHub to fetch user profiles, repositories, and README files, then generates humorous or insightful content using the Gemini AI model.
Also includes a Redis caching layer to store user profiles and repositories. Other features include roasting user resume.

[![Go](https://img.shields.io/github/go-mod/go-version/savioruz/roastgithub-api)](https://golang.org/)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/savioruz/roastgithub-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/savioruz/roastgithub-api)](https://goreportcard.com/report/github.com/savioruz/roastgithub-api)
[![License](https://img.shields.io/github/license/savioruz/roastgithub-api)]

## Table of Contents
- [Features](#features)
- [Deployment](#deployment)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
  - [Running the API](#running-the-api)
    - [Docker](#docker)
    - [Make](#make)
  - [API Documentation](#api-documentation)
- [Project Structure](#project-structure)
- [Contributing](#contributing)
- [License](#license)
- [Reference](#reference)
- [Acknowledgements](#acknowledgements)

## Features

- Fetch GitHub user profiles and repositories.
- Auto-detect language based on user location or manual input (currently id or en).
- Generate content using Google’s Gemini AI model.

## Deployment

- ### Koyeb
[![Deploy to Koyeb](https://www.koyeb.com/static/images/deploy/button.svg)](https://app.koyeb.com/services/deploy?type=git&builder=dockerfile&repository=github.com/savioruz/roastgithub-api&branch=main&ports=3000;http;/&name=roastgithub-api-koyeb&env[STAGE_STATUS]=prod&env[APP_NAME]=roastgithub-api&env[APP_HOST]=0.0.0.0&env[APP_PORT]=3000&env[GEMINI_API_KEY]=YOUR_API_KEY&env[GITHUB_TOKEN]=YOUR_GITHUB_TOKEN&env[REDIS_HOST]=YOUR_REDIS_HOST&env[REDIS_PORT]=6379&env[REDIS_PASSWORD]=&env[REDIS_DB_NUMBER]=0)

- ### Railway
[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template/jT1IvF?referralCode=XVMtOY)

- ### Render
[![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy?repo=https://github.com/savioruz/roastgithub-api)

## Requirements

- Go 1.22+
- Docker
- Redis
- Make

## Installation

1. **Clone the repository:**

    ```bash
    git clone https://github.com/savioruz/roastgithub-api.git
    cd roastgithub-api
    ```

2. **Environment Variables:**

   Create a `.env` file in the root directory and add the following:

    ```bash
    cp .env.example .env
    ```

## Usage

### Running the API

You can run the API using Docker or directly with Make.

### Docker

1. **Run redis:**

    ```bash
   make docker.redis
   ```

2. **Run the application:**

    ```bash
    make docker.run
    ```

For production, you need to secure redis on Makefile with a password.

### Make

1. **Run the application:**

    ```bash
    make run
    ```

You need to have Redis running on your machine.

### API Documentation

Swagger documentation is available at: http://localhost:3000/swagger.

## Project Structure

```
.
├── app/
│   ├── handlers/        # HTTP handler functions
│   └── models/          # Data structures and models
├── docs/                # Documentation files
├── pkg/
│   ├── middleware/      # Middleware functions for request handling
│   ├── repository/      # Data access layer
│   ├── routes/          # API route definitions
│   └── utils/           # Utility functions including GitHub and Gemini services
└── platform/
    └── cache/           # Redis caching implementation

```

## Contributing

Feel free to open issues or submit pull requests with improvements.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Reference

Inspired by the [roastgithub](https://github.com/bagusindrayana/roastgithub-api) project.

## Acknowledgements

- [Fiber](https://github.com/gofiber/fiber)
- [Go-Gemini](https://github.com/google/generative-ai-go)
- [Go-GitHub](https://github.com/google/go-github)
- [fiber-go-template](https://github.com/create-go-app/fiber-go-template)
