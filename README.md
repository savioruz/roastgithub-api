# RoastGitHub API

RoastGitHub API is a service that interacts with GitHub to fetch user profiles, repositories, and README files, then generates humorous or insightful content using the Gemini AI model.

## Features

- Fetch GitHub user profiles and repositories.
- Auto-detect language based on user location or manual input.
- Generate content using Google’s Gemini AI model.

## Requirements

- Go 1.22+
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

### Make

1. **Run the application:**

    ```bash
    make run
    ```

2. Swagger documentation is available at `http://localhost:3000/swagger`.

## Project Structure

```
.
├── app/
│   ├── handlers/        # HTTP handler functions
│   └── models/          # Data structures and models
├── docs/                # Documentation files
└── pkg/
    ├── middleware/      # Middleware functions for request handling
    ├── repository/      # Data access layer
    ├── routes/          # API route definitions
    └── utils/           # Utility functions including GitHub and Gemini services

```

## Contributing

Feel free to open issues or submit pull requests with improvements.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
