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
gcloud run services describe hema-lessons-api --region us-east1

# Get the service URL
gcloud run services describe hema-lessons-api --region us-east1 --format='value(status.url)'

# Test the health endpoint
curl "$(gcloud run services describe hema-lessons-api --region us-east1 --format='value(status.url)')/healthz"

# View logs
gcloud run services logs read hema-lessons-api --region us-east1
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
gcloud run services describe hema-lessons-api --region us-east1

# View real-time logs
gcloud run services logs tail hema-lessons-api --region us-east1

# View recent logs
gcloud run services logs read hema-lessons-api --region us-east1 --limit 100

# Update environment variables
gcloud run services update hema-lessons-api \
  --region us-east1 \
  --update-env-vars KEY=value

# Redeploy from source
gcloud run deploy hema-lessons-api --source . --region us-east1

# View revisions (for rollback)
gcloud run revisions list --service hema-lessons-api --region us-east1

# Rollback to a previous revision
gcloud run services update-traffic hema-lessons-api \
  --region us-east1 \
  --to-revisions <revision-name>=100

#deploy new code:
gcloud run deploy hema-lessons-api --source . --region us-east1
```

## Free Tier Limits

| Service | Free Tier |
|---------|-----------|
| Google Cloud Run | 2M requests, 180K vCPU-seconds, 360K GiB-seconds/month |
| Sentry | 5,000 errors/month |
| Better Stack | 5 monitors, 3-min intervals |

**Total monthly cost: $0** (within free tier limits)

> **Important**: The Cloud Run free tier applies only in Tier 1 regions (e.g., `us-east1`). Make sure to deploy to a Tier 1 region to stay within the free allowance.

---

## Mobile App Deployment

The mobile app is a bare React Native project. Builds are produced locally using Gradle (Android) and Xcode (iOS) — no Expo or EAS required.

Use the helper script from the project root:

```bash
./scripts/mobile-dev.sh <command>
```

### Android — Debug APK

A debug APK can be built on any platform with Android Studio / the Android SDK installed.

```bash
./scripts/mobile-dev.sh apk-debug
```

Output: `mobile/android/app/build/outputs/apk/debug/app-debug.apk`

Install directly on a connected device or running emulator:

```bash
adb install mobile/android/app/build/outputs/apk/debug/app-debug.apk
```

### Android — Release APK

#### 1. Generate a release keystore (first time only)

```bash
keytool -genkeypair -v \
  -keystore hema-lessons-release.keystore \
  -alias hema-lessons \
  -keyalg RSA -keysize 2048 -validity 10000
```

Store the resulting `.keystore` file somewhere safe outside the repository.

#### 2. Configure signing in `~/.gradle/gradle.properties`

```properties
HEMA_UPLOAD_STORE_FILE=/absolute/path/to/hema-lessons-release.keystore
HEMA_UPLOAD_STORE_PASSWORD=your_store_password
HEMA_UPLOAD_KEY_ALIAS=hema-lessons
HEMA_UPLOAD_KEY_PASSWORD=your_key_password
```

> **Note**: Keeping credentials in `~/.gradle/gradle.properties` (outside the repo) prevents them from being accidentally committed.

#### 3. Build

```bash
./scripts/mobile-dev.sh apk-release
```

Output: `mobile/android/app/build/outputs/apk/release/app-release.apk`

If the signing config is absent the script will warn and fall back to the debug key — the APK will work on devices but cannot be uploaded to the Play Store.

#### 4. Upload to Google Play

Upload the `.apk` (or produce an `.aab` with `./gradlew bundleRelease`) through the [Google Play Console](https://play.google.com/console).

---

### iOS — IPA for Distribution (macOS only)

> **Prerequisites**
> - macOS with Xcode installed
> - Active [Apple Developer Program](https://developer.apple.com/programs/) membership
> - CocoaPods installed (`sudo gem install cocoapods`)

#### 1. Configure `ExportOptions.plist`

A template is provided at `mobile/ios/ExportOptions.plist`. Edit it before building:

```xml
<key>teamID</key>
<string>YOUR_TEAM_ID</string>   <!-- 10-char ID from developer.apple.com/account → Membership -->

<key>method</key>
<string>app-store</string>      <!-- app-store | ad-hoc | development | enterprise -->
```

| Method | Use case |
|--------|----------|
| `app-store` | TestFlight and App Store submission |
| `ad-hoc` | Direct distribution to registered devices |
| `development` | Internal testing on registered devices |
| `enterprise` | In-house distribution (Enterprise account required) |

#### 2. Set up code signing in Xcode (first time only)

Open `mobile/ios/HemaLessons.xcworkspace` in Xcode, select the **HemaLessons** target, and under **Signing & Capabilities** choose your team and let Xcode manage signing automatically.

#### 3. Build and export the IPA

```bash
./scripts/mobile-dev.sh ios-deploy
```

The script will:
1. Run `pod install` to install CocoaPods dependencies
2. Archive the app with `xcodebuild archive`
3. Export the `.ipa` with `xcodebuild -exportArchive`

Output: `mobile/ios/build/HemaLessons.ipa`

#### 4. Upload to TestFlight / App Store

Use [Transporter](https://apps.apple.com/app/transporter/id1450874784) (free on the Mac App Store) or the Xcode Organizer to upload the `.ipa` to App Store Connect.
