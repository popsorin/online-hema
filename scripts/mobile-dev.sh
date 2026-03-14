#!/bin/bash
# React Native Mobile Development Helper Script
#
# Usage:
#   ./scripts/mobile-dev.sh android      - Start emulator if needed, then run app on it
#   ./scripts/mobile-dev.sh ios          - Run app on iOS simulator (macOS only)
#   ./scripts/mobile-dev.sh apk-debug    - Build a debug APK
#   ./scripts/mobile-dev.sh apk-release  - Build a release APK
#   ./scripts/mobile-dev.sh ios-deploy   - Archive and export an IPA for distribution (macOS only)
#   ./scripts/mobile-dev.sh start        - Start Metro bundler only
#   ./scripts/mobile-dev.sh test         - Run mobile tests
#   ./scripts/mobile-dev.sh lint         - Run linting and type checking
#   ./scripts/mobile-dev.sh install      - Install/update dependencies
#   ./scripts/mobile-dev.sh clean        - Clean and rebuild
#
# Prerequisites:
#   - Android Studio with an AVD configured (for Android)
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
        RUNNING_EMULATOR=$(adb devices 2>/dev/null | grep -E "^emulator-" | awk '{print $1}')
        if [ -z "$RUNNING_EMULATOR" ]; then
            AVD_NAME=$(emulator -list-avds 2>/dev/null | head -1)
            if [ -z "$AVD_NAME" ]; then
                log_error "No Android Virtual Devices found. Create one in Android Studio."
                exit 1
            fi
            log_info "No emulator running. Starting AVD: $AVD_NAME"
            emulator -avd "$AVD_NAME" -no-audio -no-snapshot-save &>/dev/null &
            log_info "Waiting for emulator to boot..."
            adb wait-for-device
            until adb shell getprop sys.boot_completed 2>/dev/null | grep -q "1"; do
                sleep 3
            done
            log_info "Emulator ready."
        else
            log_info "Emulator already running: $RUNNING_EMULATOR"
        fi
        echo ""
        npx react-native run-android
        ;;

    ios)
        log_info "Starting app on iOS simulator..."
        if [[ "$OSTYPE" != "darwin"* ]]; then
            log_error "iOS development is only available on macOS"
            exit 1
        fi
        npx react-native run-ios
        ;;

    apk-debug)
        log_info "Building debug APK..."
        cd "$MOBILE_DIR/android"
        ./gradlew assembleDebug
        APK_PATH="$MOBILE_DIR/android/app/build/outputs/apk/debug/app-debug.apk"
        if [ -f "$APK_PATH" ]; then
            log_info "Debug APK built successfully:"
            log_info "  $APK_PATH"
            echo ""
            log_info "To install on a connected device or emulator:"
            log_info "  adb install \"$APK_PATH\""
        fi
        ;;

    apk-release)
        log_info "Building release APK..."
        if ! grep -qs "HEMA_UPLOAD_STORE_FILE" "$HOME/.gradle/gradle.properties" 2>/dev/null; then
            log_warn "No release signing config found in ~/.gradle/gradle.properties"
            log_warn "The APK will be signed with the debug key (not suitable for Play Store)"
            log_warn "See docs/deployment.md for signing setup instructions"
            echo ""
        fi
        cd "$MOBILE_DIR/android"
        ./gradlew assembleRelease
        APK_PATH="$MOBILE_DIR/android/app/build/outputs/apk/release/app-release.apk"
        if [ -f "$APK_PATH" ]; then
            log_info "Release APK built successfully:"
            log_info "  $APK_PATH"
        fi
        ;;

    ios-deploy)
        if [[ "$OSTYPE" != "darwin"* ]]; then
            log_error "iOS deployment is only available on macOS"
            exit 1
        fi
        EXPORT_OPTIONS="$MOBILE_DIR/ios/ExportOptions.plist"
        if [ ! -f "$EXPORT_OPTIONS" ]; then
            log_error "ExportOptions.plist not found at: $EXPORT_OPTIONS"
            log_error "See docs/deployment.md for instructions on creating it"
            exit 1
        fi
        WORKSPACE="$MOBILE_DIR/ios/HemaLessons.xcworkspace"
        ARCHIVE_PATH="$MOBILE_DIR/ios/HemaLessons.xcarchive"
        EXPORT_PATH="$MOBILE_DIR/ios/build"
        log_info "Installing CocoaPods dependencies..."
        cd "$MOBILE_DIR/ios" && pod install
        log_info "Archiving the app (this may take a few minutes)..."
        xcodebuild archive \
            -workspace "$WORKSPACE" \
            -scheme HemaLessons \
            -configuration Release \
            -archivePath "$ARCHIVE_PATH"
        log_info "Exporting IPA..."
        xcodebuild -exportArchive \
            -archivePath "$ARCHIVE_PATH" \
            -exportOptionsPlist "$EXPORT_OPTIONS" \
            -exportPath "$EXPORT_PATH"
        IPA_PATH="$EXPORT_PATH/HemaLessons.ipa"
        if [ -f "$IPA_PATH" ]; then
            log_info "IPA exported successfully:"
            log_info "  $IPA_PATH"
            echo ""
            log_info "Upload to TestFlight / App Store via Transporter or Xcode Organizer"
        fi
        ;;

    start)
        log_info "Starting Metro bundler..."
        npx react-native start
        ;;

    test)
        log_info "Running mobile tests..."
        npm test
        ;;

    lint)
        log_info "Running lint and type checks..."
        npm run lint && npm run typecheck
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
        echo "  android      Start emulator if needed, then run app on it"
        echo "  ios          Run app on iOS simulator (macOS only)"
        echo "  apk-debug    Build a debug APK"
        echo "  apk-release  Build a signed release APK"
        echo "  ios-deploy   Archive and export an IPA for distribution (macOS only)"
        echo "  start        Start Metro bundler only"
        echo "  test         Run mobile unit tests"
        echo "  lint         Run ESLint and TypeScript checks"
        echo "  install      Install npm dependencies"
        echo "  clean        Clean builds and reinstall dependencies"
        echo "  help         Show this help message"
        echo ""
        echo "Prerequisites:"
        echo "  - Node.js 18+"
        echo "  - Android Studio with an AVD configured (for Android)"
        echo "  - Xcode (for iOS, macOS only)"
        echo ""
        echo "Examples:"
        echo "  $0 android       # Run on Android emulator"
        echo "  $0 apk-debug     # Build a debug APK"
        echo "  $0 apk-release   # Build a release APK"
        echo "  $0 ios-deploy    # Archive and export IPA (macOS)"
        echo "  $0 test          # Run tests locally"
        ;;
esac
