#!/usr/bin/env bash

if [ -z "$MM_WEBHOOK_URL" ]
then
  echo "Please set MM_WEBHOOK_URL environment variable"
  exit 0
fi

LOCAL_PLUGINSJSON_PATH="plugins.json"
REMOTE_PLUGINSJSON_PATH="new_plugins.json"

curl https://raw.githubusercontent.com/mattermost/mattermost-marketplace/production/plugins.json -o $REMOTE_PLUGINSJSON_PATH

CHANGES="$(git diff --no-index $LOCAL_PLUGINSJSON_PATH $REMOTE_PLUGINSJSON_PATH)"

# Save the original plugins.json as a backup
mv $LOCAL_PLUGINSJSON_PATH ${LOCAL_PLUGINSJSON_PATH}.bak
mv $REMOTE_PLUGINSJSON_PATH $LOCAL_PLUGINSJSON_PATH

if [ -z "$CHANGES" ]
then
  echo "no changes to plugins.json"
  exit 0
fi

JSON="$(echo "$CHANGES" | grep '^+ ' | sed 's/^+ //g' | sed '$s/},/}/g')"

NAME=$(jq -r .manifest.name <<< "$JSON")
VERSION=$(jq -r .manifest.version <<< "$JSON")
MIN_SERVER_VERSION=$(jq -r .manifest.min_server_version <<< "$JSON")
RELEASE_NOTES=$(jq -r .release_notes_url <<< "$JSON")

echo '{
  "username": "Plugin Marketplace",
  "icon_url": "https://mattermost.com/wp-content/uploads/2022/02/icon.png",
  "attachments":[{
    "fallback": "'$NAME $VERSION' was added to the Marketplace",
    "title": "'$NAME $VERSION' was added to the Marketplace",
    "text": "Release notes can be found [here]('$RELEASE_NOTES'). It requires Mattermost Server '$MIN_SERVER_VERSION'."
  }]
}' > mattermost.json

curl -i -X POST -H 'Content-Type: application/json' -d @mattermost.json $MM_WEBHOOK_URL
cat mattermost.json
