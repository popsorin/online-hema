#!/bin/bash
# React Native Mobile Development Helper Script
#
# Usage:
#   ./scripts/mobile-dev.sh android   - Run app on Android emulator
#   ./scripts/mobile-dev.sh ios       - Run app on iOS simulator  
#   ./scripts/mobile-dev.sh start     - Start Metro bundler only
#   ./scripts/mobile-dev.sh test      - Run mobile tests
#   ./scripts/mobile-dev.sh lint      - Run linting and type checking
#   ./scripts/mobile-dev.sh install   - Install/update dependencies
#   ./scripts/mobile-dev.sh clean     - Clean and rebuild
#
# Prerequisites:
#   - Android Studio with emulator (for Android)
#   - Xcode (for iOS, macOS only)

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
MOBILE_DIR="$PROJECT_ROOT/mobile"

cd "$MOBILE_DIR"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

case "${1:-help}" in
    android)
        log_info "Starting app on Android emulator..."
        log_info "Make sure you have an Android emulator running!"
        echo ""
        npx react-native run-android
        ;;
        
    restart)
        log_info "Restarting Android app..."
        adb shell am force-stop com.hemalessons
        sleep 3
        adb shell am start -n com.hemalessons/.MainActivity
        log_info "App restarted!"
        ;;
        
    ios)
        log_info "Starting app on iOS simulator..."
        if [[ "$OSTYPE" != "darwin"* ]]; then
            log_error "iOS development is only available on macOS"
            exit 1
        fi
        npx react-native run-ios
        ;;
        
    start)
        log_info "Starting Metro bundler..."
        npx react-native start
        ;;
        
    test)
        log_info "Running mobile tests..."
        npm test
        ;;
        
    test-docker)
        log_info "Running mobile tests in Docker..."
        cd "$PROJECT_ROOT"
        docker compose -f docker-compose.mobile.yml --profile test run --rm mobile-test
        ;;
        
    lint)
        log_info "Running lint and type checks..."
        npm run lint && npm run typecheck
        ;;
        
    lint-docker)
        log_info "Running lint and type checks in Docker..."
        cd "$PROJECT_ROOT"
        docker compose -f docker-compose.mobile.yml --profile ci run --rm mobile-ci
        ;;
        
    install)
        log_info "Installing mobile dependencies..."
        npm install
        ;;
        
    clean)
        log_info "Cleaning build caches..."
        rm -rf node_modules
        rm -rf android/build android/.gradle
        rm -rf ios/build ios/Pods
        npm install
        log_info "Clean complete."
        ;;
        
    help|*)
        echo "HEMA Lessons React Native Mobile Development Helper"
        echo ""
        echo "Usage: $0 <command>"
        echo ""
        echo "Commands:"
        echo "  android      Run app on Android emulator"
        echo "  ios          Run app on iOS simulator (macOS only)"
        echo "  start        Start Metro bundler only"
        echo "  restart      Force stop and restart Android app"
        echo "  test         Run mobile unit tests locally"
        echo "  test-docker  Run mobile unit tests in Docker"
        echo "  lint         Run ESLint and TypeScript checks locally"
        echo "  lint-docker  Run lint checks in Docker"
        echo "  install      Install npm dependencies"
        echo "  clean        Clean builds and reinstall dependencies"
        echo "  help         Show this help message"
        echo ""
        echo "Prerequisites:"
        echo "  - Node.js 18+"
        echo "  - Android Studio with emulator (for Android)"
        echo "  - Xcode (for iOS, macOS only)"
        echo ""
        echo "Examples:"
        echo "  $0 android    # Run on Android emulator"
        echo "  $0 test       # Run tests locally"
        ;;
esac
