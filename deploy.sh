#!/bin/bash

# Stop execution if any command fails
set -e

echo "Starting Deployment Process..."
echo "---------------------------------"

if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
    echo ".env loaded."
else
    echo "Error: .env file not found."
    exit 1
fi

if [ -z "$DATABASE_URL" ]; then
    echo "Error: DATABASE_URL is empty. Check your .env file."
    exit 1
fi

if [ -z "$PERSPECTIVE_API_KEY" ]; then
    echo "Error: PERSPECTIVE_API_KEY is empty. Check your .env file."
    exit 1
fi

# Build the new container image using Dockerfile
echo "Building Container..."
gcloud builds submit --tag gcr.io/portfolio-site-480019/portfolio

# Deploy the new image to Cloud Run
echo "Deploying to Cloud Run (us-east1)..."
gcloud run deploy portfolio \
  --image gcr.io/portfolio-site-480019/portfolio \
  --platform managed \
  --region us-east1 \
  --allow-unauthenticated \
  --port 8080 \
  --memory 512Mi \
  --cpu-boost \
  --set-env-vars DATABASE_URL="$DATABASE_URL",PERSPECTIVE_API_KEY="$PERSPECTIVE_API_KEY"

echo "---------------------------------"
echo "✅ Deployment Complete!"