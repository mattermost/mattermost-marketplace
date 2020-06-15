#!/usr/bin/env bash

set -euxo pipefail
CHANGES=$(git diff HEAD~1 plugins.json)
if [ -z "$CHANGES" ]
then
  echo "no changes to plugins.json"
  echo "{}" > mattermost.json
  exit 0
fi

JSON=$(git diff HEAD~1 plugins.json | grep '^+ ' | sed 's/+//g' | sed '$s/},/}/g')
NAME=$(jq -r .manifest.name <<< "$JSON")
VERSION=$(jq -r .manifest.version <<< "$JSON")
MIN_SERVER_VERSION=$(jq -r .manifest.min_server_version <<< "$JSON")
RELEASE_NOTES=$(jq -r .release_notes_url <<< "$JSON")

echo '{
  "username": "Plugin Marketplace",
  "icon_url": "https://www.mattermost.org/wp-content/uploads/2016/04/icon.png",
  "attachments":[{
    "fallback": "'$NAME $VERSION' was added to the Marketplace",
    "title": "'$NAME $VERSION' was added to the Marketplace",
    "text": "Release notes can be found [here]('$RELEASE_NOTES'). It requires Mattermost Server '$MIN_SERVER_VERSION'."
  }]
}' > mattermost.json
cat mattermost.json
