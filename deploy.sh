#!/bin/bash

# Stop execution if any command fails
set -e

echo "Starting Deployment Process..."
echo "---------------------------------"

# Build the new container image using Dockerfile
echo "Building Container..."
gcloud builds submit --tag gcr.io/portfolio-site-480019/portfolio

# Deploy the new image to Cloud Run
echo "☁️  Deploying to Cloud Run (us-east1)..."
gcloud run deploy portfolio \
  --image gcr.io/portfolio-site-480019/portfolio \
  --platform managed \
  --region us-east1 \
  --allow-unauthenticated \
  --port 8080

echo "---------------------------------"
echo "✅ Deployment Complete!"