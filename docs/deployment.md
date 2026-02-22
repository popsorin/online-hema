# Deployment Guide

This guide covers deploying the HEMA Lessons API to Google Cloud Run and setting up monitoring.

## Prerequisites

- [Google Cloud account](https://cloud.google.com/) (free tier)
- [gcloud CLI](https://cloud.google.com/sdk/docs/install) installed
- [Sentry account](https://sentry.io/signup/) (free tier)
- [Better Stack account](https://betterstack.com/uptime) (free tier)

## 1. Deploy to Google Cloud Run

### Install the gcloud CLI

```bash
# macOS
brew install --cask google-cloud-sdk

# Linux
curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-cli-linux-x86_64.tar.gz
tar -xf google-cloud-cli-linux-x86_64.tar.gz
./google-cloud-sdk/install.sh

# Windows
# Download installer from https://cloud.google.com/sdk/docs/install
```

### Authenticate

```bash
gcloud auth login
```

### First-time Setup

```bash
# Create a new project (or use an existing one)
gcloud projects create hema-lessons-api --name="HEMA Lessons API"
gcloud config set project hema-lessons-api

# Enable required APIs
gcloud services enable \
  run.googleapis.com \
  artifactregistry.googleapis.com \
  cloudbuild.googleapis.com

# Link a billing account (required even for free tier)
# Visit https://console.cloud.google.com/billing to set up billing
```

### Deploy

The `--source .` flag uses Cloud Build to build the existing Dockerfile and push the image to Artifact Registry automatically.

```bash
gcloud run deploy hema-lessons-api \
  --source . \
  --region us-east1 \
  --allow-unauthenticated \
  --port 8080 \
  --memory 256Mi \
  --cpu 1 \
  --min-instances 0 \
  --max-instances 2 \
  --set-env-vars APP_ENVIRONMENT=production
```

When prompted, confirm the Artifact Registry repository creation.

### Verify Deployment

```bash
# Check service status
gcloud run services describe hema-lessons-api --region us-central1

# Get the service URL
gcloud run services describe hema-lessons-api --region us-central1 --format='value(status.url)'

# Test the health endpoint
curl "$(gcloud run services describe hema-lessons-api --region us-central1 --format='value(status.url)')/healthz"

# View logs
gcloud run services logs read hema-lessons-api --region us-central1
```

Your API will be available at: `https://hema-lessons-api-<hash>.run.app`

## 2. Set Up Sentry (Error Tracking)

1. Create a new project at [sentry.io](https://sentry.io)
   - Platform: Go
   - Project name: hema-lessons-api

2. Copy the DSN from project settings

3. Set the Sentry DSN as an environment variable on Cloud Run:
   ```bash
   gcloud run services update hema-lessons-api \
     --region us-east1 \
     --update-env-vars SENTRY_DSN=https://your-key@sentry.io/project-id
   ```

## 3. Set Up Better Stack (Uptime Monitoring)

1. Sign up at [betterstack.com/uptime](https://betterstack.com/uptime)

2. Create a new monitor:
   - **URL**: `https://hema-lessons-api-<hash>.run.app/healthz`
   - **Check interval**: 3 minutes (free tier)
   - **Request timeout**: 30 seconds
   - **HTTP Method**: GET

3. Configure alerts:
   - Email notifications (free)
   - Optional: Slack/SMS integration

4. (Optional) Create a status page:
   - Add your monitor to a public status page
   - Share URL with users

> **Note**: Better Stack health checks will also help keep the Cloud Run service warm, reducing cold starts.

## Useful Commands

```bash
# View service status
gcloud run services describe hema-lessons-api --region us-central1

# View real-time logs
gcloud run services logs tail hema-lessons-api --region us-central1

# View recent logs
gcloud run services logs read hema-lessons-api --region us-central1 --limit 100

# Update environment variables
gcloud run services update hema-lessons-api \
  --region us-central1 \
  --update-env-vars KEY=value

# Redeploy from source
gcloud run deploy hema-lessons-api --source . --region us-central1

# View revisions (for rollback)
gcloud run revisions list --service hema-lessons-api --region us-central1

# Rollback to a previous revision
gcloud run services update-traffic hema-lessons-api \
  --region us-central1 \
  --to-revisions <revision-name>=100
```

## Free Tier Limits

| Service | Free Tier |
|---------|-----------|
| Google Cloud Run | 2M requests, 180K vCPU-seconds, 360K GiB-seconds/month |
| Sentry | 5,000 errors/month |
| Better Stack | 5 monitors, 3-min intervals |

**Total monthly cost: $0** (within free tier limits)

> **Important**: The Cloud Run free tier applies only in Tier 1 regions (e.g., `us-central1`). Make sure to deploy to a Tier 1 region to stay within the free allowance.
