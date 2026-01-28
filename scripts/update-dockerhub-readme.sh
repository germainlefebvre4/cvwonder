#!/bin/bash
# Script to automatically update Docker tags in docs/dockerhub/README.md

set -e

DOCKERHUB_README="docs/dockerhub/README.md"
CHANGELOG="CHANGELOG.md"

# Extract the latest version from CHANGELOG.md
# Look for the first version heading like ## [0.2.1] or ## [0.2.0]
LATEST_VERSION=$(grep -oP '## \[\K[0-9]+\.[0-9]+\.[0-9]+' "$CHANGELOG" | head -n 1)

if [ -z "$LATEST_VERSION" ]; then
    echo "Error: Could not extract version from CHANGELOG.md"
    exit 1
fi

# Build the version tags
MAJOR=$(echo "$LATEST_VERSION" | cut -d. -f1)
MINOR=$(echo "$LATEST_VERSION" | cut -d. -f2)
PATCH=$(echo "$LATEST_VERSION" | cut -d. -f3)

# Create tags like: latest, v0, v0.2, v0.2.1
TAGS="\`latest\`, \`v${MAJOR}\`, \`v${MAJOR}.${MINOR}\`, \`v${MAJOR}.${MINOR}.${PATCH}\`"
TAG_LINE="* [$TAGS](https://github.com/germainlefebvre4/kuve/blob/v${LATEST_VERSION}/Dockerfile)"

# Update the DockerHub README
# Replace the line that starts with "* [\`latest\`," with the new tag line
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    sed -i '' "s|^\* \[\`latest\`.*|$TAG_LINE|" "$DOCKERHUB_README"
else
    # Linux
    sed -i "s|^\* \[\`latest\`.*|$TAG_LINE|" "$DOCKERHUB_README"
fi

echo "Updated $DOCKERHUB_README with version v${LATEST_VERSION}"
echo "Tags: $TAGS"
