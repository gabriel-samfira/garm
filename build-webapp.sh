#!/bin/bash

set -e

echo "Building GARM SPA (SvelteKit)..."

# Navigate to webapp directory
cd webapp

# Install dependencies if node_modules doesn't exist
if [ ! -d "node_modules" ]; then
    echo "Installing dependencies..."
    npm install
fi

# Build the SPA
echo "Building SPA..."
npm run build
echo "SPA built successfully!"
