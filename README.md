# New Project

A simple Go web application set up for deployment with Google Cloud Build and Cloud Run.

## Local Development

To run the application locally:

```bash
go run cmd/new_project/main.go
```

The application will be available at http://localhost:8080

## Building the Docker Image

```bash
docker build -t new-project .
docker run -p 8080:8080 new-project
```

## Deployment

This project is set up for automatic deployment using Cloud Build and Cloud Run.

### Prerequisites

1. Enable the following APIs in your Google Cloud project:
   - Cloud Build API
   - Cloud Run API
   - Container Registry API

2. Set up a Cloud Build trigger for your GitHub repository.

### Manual Deployment

You can also deploy manually using gcloud:

```bash
gcloud builds submit --config cloudbuild.yaml
```