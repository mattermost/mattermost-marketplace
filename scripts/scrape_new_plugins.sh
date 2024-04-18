#!/usr/bin/env bash

if [ -z "$MM_WEBHOOK_URL" ]
then
  echo "Please set the MM_WEBHOOK_URL environment variable"
  exit 0
fi

LOCAL_PLUGINSJSON_PATH="plugins.json"
REMOTE_PLUGINSJSON_PATH="new_plugins.json"

REMOTE_URL=https://raw.githubusercontent.com/mattermost/mattermost-marketplace/production/plugins.json
echo ""
echo "Fetching remote plugins from $REMOTE_URL"
echo ""

curl $REMOTE_URL -o $REMOTE_PLUGINSJSON_PATH --silent

if [ ! -f $LOCAL_PLUGINSJSON_PATH ]; then
  cp $REMOTE_PLUGINSJSON_PATH $LOCAL_PLUGINSJSON_PATH
  echo "Saved remote marketplace entries to plugin.json"
  echo "First time running program"
  echo "Exiting with nothing to compare"
  exit 0
fi

local_keys=$(jq -r '.[] | .manifest.id + "-" + .manifest.version' $LOCAL_PLUGINSJSON_PATH)
remote_keys=$(jq -r '.[] | .manifest.id + "-" + .manifest.version' $REMOTE_PLUGINSJSON_PATH)

# Backup plugins.json file
mv $LOCAL_PLUGINSJSON_PATH ${LOCAL_PLUGINSJSON_PATH}.bak
mv $REMOTE_PLUGINSJSON_PATH $LOCAL_PLUGINSJSON_PATH

local_array=($local_keys)
remote_array=($remote_keys)

HAS_NEW_PLUGINS=0
for rkey in "${remote_array[@]}"; do
    FOUND_EXISTING_ENTRY=0
    for lkey in "${local_array[@]}"; do
        if [ "$rkey" == "$lkey" ]; then
            FOUND_EXISTING_ENTRY=1
            break
        fi
    done
    if [ $FOUND_EXISTING_ENTRY -ne 1 ]; then
        if [ $HAS_NEW_PLUGINS -ne 1 ]; then
          echo "New plugins found:"
          HAS_NEW_PLUGINS=1
        fi

        # Split the key into id and version. This handles the case where the plugin's id has dashes in it.
        ID="${rkey%-*}"
        VERSION="${rkey##*-}"
        echo "- $ID v$VERSION"

        # Extract the full object that matches the id and version
        JSON=$(jq --arg id "$ID" --arg version "$VERSION" \
            '.[] | select(.manifest.id == $id and .manifest.version == $version)' $LOCAL_PLUGINSJSON_PATH)

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

        # Send webhook request to Mattermost
        curl -i -X POST -H 'Content-Type: application/json' -d @mattermost.json $MM_WEBHOOK_URL --silent --output /dev/null --show-error --fail
    fi
done

if [ $HAS_NEW_PLUGINS -ne 1 ]; then
  echo "No new plugins found when comparing remote marketplace to local copy"
fi
