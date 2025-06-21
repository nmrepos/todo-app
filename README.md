# Todo Application 

A simple and intuitive web-based Todo application built with Go and Gin for the backend and vanilla JavaScript and CSS for the frontend. The project includes a robust CI/CD pipeline utilizing GitHub Actions to ensure quality and easy deployment.

## Features

- **Create Tasks:** Easily add new todo items.
- **Read Tasks:** View all existing todo items.
- **Update Tasks:** Mark tasks as completed or update their titles.
- **Delete Tasks:** Remove tasks no longer needed.

## Project Structure

```
.
├── cmd
│   └── server
│       ├── main.go
│       ├── main_test.go
│       ├── complete_test.go
│       └── delete_test.go
├── static
│   └── index.html
├── .github
│   └── workflows
│       └── ci.yml
├── Dockerfile
├── go.mod
├── go.sum
├── .gitignore
└── README.md
```

## Setup Instructions

### Prerequisites

- Go (version 1.20 or higher)
- Docker
- Git

### Clone the Repository

```bash
git clone https://github.com/yourusername/todo-app.git
cd todo-app
```

### Running Locally

**Install Dependencies:**

```bash
go mod download
```

**Run the Server:**

```bash
go run ./cmd/server/main.go
```

Open your browser and go to: `http://localhost:8080`

### Docker Setup

Build Docker image:

```bash
docker build -t todo-app .
```

Run Docker container:

```bash
docker run -d -p 8080:8080 todo-app
```

Access the application via: `http://localhost:8080`

## Testing

Run unit and integration tests with:

```bash
go test ./cmd/server -v
```

## CI/CD Pipeline

This project uses GitHub Actions for continuous integration and deployment.

### CI Pipeline Overview

- **Linting:** Ensures code quality with `golangci-lint`.
- **Build:** Verifies the application compiles correctly.
- **Test:** Runs unit and integration tests, ensuring all endpoints work correctly.
- **Docker Image:** Builds and pushes Docker images to Docker Hub.
- **Deployment:**
  - **Dev Environment:** Automatically deploys on push to `develop`.
  - **Prod Environment:** Manually triggered or automatically upon merging into `main`.

### Triggering CI/CD

- Push commits to `develop` to trigger automatic deployment to development.
- Merge into `main` or manually dispatch workflow for production deployment.

### Setting up Docker Hub Credentials for CI/CD

- To allow GitHub Actions to push Docker images to your Docker Hub account:
- Generate Docker Hub Access Token
- Log in to Docker Hub.
- Go to Account Settings > Security.
- Click New Access Token and copy the generated token.
- Add Secrets to Your GitHub Repository
- Go to your GitHub repo.
- Click Settings > Secrets and variables > Actions.
- Click New repository secret for each:
```
DOCKERHUB_USERNAME — Your Docker Hub username.

DOCKERHUB_TOKEN — The access token you generated.
```

### Docker Hub Images

Docker images are pushed to Docker Hub with the following tags:

- `dev`: Latest development version.
- `latest`: Stable production version.

Pull an image from Docker Hub:

```bash
docker pull yourdockerhubusername/todo-app:latest
```



## Contributing

1. Fork the repository.
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Commit your changes: `git commit -m "feat(component): detailed message"`
4. Push to your feature branch: `git push origin feature/your-feature`
5. Open a Pull Request into the `develop` branch.

