# Android Build and Install Guide

This app is a React Native CLI project (`mobile/`), so Android binaries are built with Gradle from `mobile/android/`.

## 1) Prerequisites

- Android phone with **Developer options** enabled
- **USB debugging** enabled on the phone
- `adb` installed and available in PATH
- JDK installed (Java 17+ is fine for this project)

Verify:

```bash
adb version
java -version
adb devices -l
```

If `adb devices -l` shows no devices, reconnect USB, accept the RSA prompt on phone, and set USB mode to File Transfer.

## 2) Fast local install (debug)

From `mobile/android/`:

```bash
./gradlew installDebug
```

Debug APK output:

`mobile/android/app/build/outputs/apk/debug/app-debug.apk`

If install fails with `No connected devices`, connect a phone (or start an emulator) and run again.

## 3) Build release APK

From `mobile/android/`:

```bash
./gradlew assembleRelease
```

Release APK output:

`mobile/android/app/build/outputs/apk/release/app-release.apk`

Install over USB:

```bash
adb install -r /absolute/path/to/mobile/android/app/build/outputs/apk/release/app-release.apk
```

Or copy the APK to the phone and open it (allow unknown sources if prompted).

## 4) Production signing (recommended)

Right now release can fall back to debug signing when no keystore is configured. For real distribution, use your own keystore.

Create keystore (example):

```bash
keytool -genkeypair -v -storetype PKCS12 -keystore my-upload-key.keystore -alias my-key-alias -keyalg RSA -keysize 2048 -validity 10000
```

Place the keystore at `mobile/android/app/` (or another safe location), then set in `mobile/android/gradle.properties`:

```properties
HEMA_UPLOAD_STORE_FILE=my-upload-key.keystore
HEMA_UPLOAD_KEY_ALIAS=my-key-alias
HEMA_UPLOAD_STORE_PASSWORD=your-store-password
HEMA_UPLOAD_KEY_PASSWORD=your-key-password
```

Then rebuild:

```bash
cd mobile/android
./gradlew clean assembleRelease
```

## What a keystore is

A keystore is a secure file containing your app signing key. Android uses this signature to verify app updates come from the same publisher. If you lose this key, users cannot upgrade from older signed versions, so back it up securely.
