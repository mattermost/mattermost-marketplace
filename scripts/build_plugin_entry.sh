#!/bin/bash
set -euo pipefail

if [ "$#" -ne 1 ]; then
    echo "Missing DownloadURL"
    exit 1
fi

DownloadURL=$1
Repository=`echo $1 | perl -pe 's/^.+(mattermost-plugin-[a-zA-Z0-9-]+)\/.+$/\1/'`
HomepageURL=`echo $1 | perl -pe 's/(mattermost-plugin-[a-zA-Z0-9-]+)\/.+/\1/'`
DownloadSignature=
ManifestURL="https://raw.githubusercontent.com/mattermost/$Repository/master/plugin.json"
Manifest=`curl -s $ManifestURL`

echo "{\"HomepageURL\": \"$HomepageURL\", \"DownloadURL\": \"$DownloadURL\", \"DownloadSignature\": \"$DownloadSignature\", \"Manifest\": $Manifest}" | jq
