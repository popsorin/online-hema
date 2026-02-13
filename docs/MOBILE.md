# Mobile App Development Guide

This guide covers setting up and running the HEMA Lessons mobile app using Expo.

## Prerequisites

- Docker and Docker Compose installed
- Expo Go app installed on your phone (App Store / Play Store)

No need for Android Studio, Xcode, or any native SDKs - Expo handles everything.

## Quick Start

### Start Development Server

```bash
# Start full stack (API + Expo dev server)
./scripts/mobile-dev.sh start
```

This will:
1. Start the PostgreSQL database
2. Start Redis cache
3. Start the Go API server on `http://localhost:8080`
4. Start the Expo dev server

### Test on Your Phone

1. Install **Expo Go** from the App Store (iOS) or Play Store (Android)
2. Run `./scripts/mobile-dev.sh start`
3. Scan the QR code shown in the terminal
4. The app opens instantly on your phone

### Using Docker Compose Directly

```bash
# Create the shared network (first time only)
docker network create hema-lessons-network

# Start everything
docker compose -f docker-compose.yml -f docker-compose.mobile.yml up --build
```

### Start Mobile Only (API Running Separately)

If you already have the API running:

```bash
./scripts/mobile-dev.sh mobile-only
```

Or with a custom API URL:

```bash
API_URL=http://192.168.1.100:8080 ./scripts/mobile-dev.sh mobile-only
```

## Available Commands

### Helper Script Commands

```bash
./scripts/mobile-dev.sh start        # Start full stack (API + Expo)
./scripts/mobile-dev.sh mobile-only  # Start Expo dev server only
./scripts/mobile-dev.sh test         # Run unit tests
./scripts/mobile-dev.sh lint         # Run ESLint and TypeScript checks
./scripts/mobile-dev.sh install      # Install/update npm dependencies
./scripts/mobile-dev.sh clean        # Remove Docker volumes and clean cache
```

### Docker Compose Commands

```bash
# Start services
docker compose -f docker-compose.yml -f docker-compose.mobile.yml up

# Run tests
docker compose -f docker-compose.mobile.yml run --rm mobile-test

# Run linting
docker compose -f docker-compose.mobile.yml run --rm mobile-ci

# View logs
docker compose -f docker-compose.mobile.yml logs -f mobile

# Stop services
docker compose -f docker-compose.yml -f docker-compose.mobile.yml down
```

## Project Structure

```
mobile/
├── src/
│   ├── api/           # API client and endpoints
│   ├── components/    # Reusable UI components
│   ├── screens/       # Screen components
│   ├── navigation/    # React Navigation setup
│   ├── hooks/         # Custom React hooks
│   ├── store/         # State management (Zustand)
│   ├── utils/         # Utility functions and config
│   ├── types/         # TypeScript type definitions
│   └── assets/        # Images, fonts, etc.
├── Dockerfile         # Docker build configuration
├── App.tsx            # Root component
├── app.json           # Expo configuration
└── package.json       # Dependencies and scripts
```

## Configuration

### API URL

The API URL is configured in `src/utils/config.ts`. By default:

| Environment | Default URL | Notes |
|-------------|-------------|-------|
| Development (Docker) | `http://api:8080` | Service name on shared network |
| Development (Phone) | `http://10.0.2.2:8080` | Android default |
| Production | `https://api.hema-lessons.com` | Production API |

### Environment Variables

When using Docker, you can override settings via environment variables:

```bash
# Custom API URL
API_URL=http://192.168.1.100:8080 docker compose -f docker-compose.mobile.yml up mobile
```

## Expo Ports

| Port | Purpose |
|------|---------|
| 8081 | Metro bundler |
| 19000 | Expo dev server |
| 19001 | Expo dev tools |
| 19002 | Expo dev tools web UI |

## Troubleshooting

### QR Code Not Working

If scanning the QR code doesn't work:

1. Make sure your phone and computer are on the same network
2. Try using tunnel mode (default): `expo start --tunnel`
3. Check if any firewall is blocking ports 19000-19002

### Can't Connect to API from Phone

1. The API runs on `localhost:8080` which your phone can't access directly
2. Use your machine's local IP instead:
   ```bash
   API_URL=http://192.168.1.X:8080 ./scripts/mobile-dev.sh mobile-only
   ```
3. Find your IP with: `ip addr` (Linux) or `ifconfig` (macOS)

### Hot Reload Not Working

```bash
# Clean and restart
./scripts/mobile-dev.sh clean
./scripts/mobile-dev.sh start
```

### TypeScript/ESLint Errors

```bash
# Run checks to see all issues
./scripts/mobile-dev.sh lint
```

## Testing

Run the test suite:

```bash
# Via Docker
./scripts/mobile-dev.sh test

# With coverage (run inside container)
docker compose -f docker-compose.mobile.yml run --rm mobile-test npm run test:coverage
```

## Building for Production

When you're ready to publish:

1. Create an Expo account at https://expo.dev
2. Install EAS CLI: `npm install -g eas-cli`
3. Build for stores:
   ```bash
   eas build --platform android
   eas build --platform ios
   ```

See [Expo Build Documentation](https://docs.expo.dev/build/introduction/) for more details.
