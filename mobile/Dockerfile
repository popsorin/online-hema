# React Native CI/Testing Dockerfile
# This image is used for:
# - Running tests in CI/CD
# - Running lint and typecheck
#
# NOTE: Development with React Native CLI requires Android Studio/Xcode
# and should be done locally, not in Docker.

FROM node:20-alpine AS base

# Install dependencies required for builds
RUN apk add --no-cache \
    git \
    bash \
    curl

WORKDIR /app

# Copy package files first for better caching
COPY package*.json ./

# ============================================
# Test Stage
# ============================================
FROM base AS test

# Install all dependencies
RUN npm install

# Copy source code
COPY . .

# Run tests
CMD ["npm", "test"]

# ============================================
# CI/Lint Stage
# ============================================
FROM base AS ci

RUN npm install

COPY . .

# Run linting and type checking
CMD ["npm", "run", "ci:check"]
