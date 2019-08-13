#!/bin/bash
set -euo pipefail

DIRECTORY=$(cd `dirname $0` && pwd)

echo '['
$DIRECTORY/build_plugin_entry.sh https://github.com/mattermost/mattermost-plugin-demo/releases/download/v0.1.0/com.mattermost.demo-plugin-0.1.0.tar.gz
echo ','
$DIRECTORY/build_plugin_entry.sh https://github.com/mattermost/mattermost-plugin-github/releases/download/v0.10.2/github-0.10.2.tar.gz
echo ','
$DIRECTORY/build_plugin_entry.sh https://github.com/mattermost/mattermost-plugin-autolink/releases/download/v1.0.1/mattermost-autolink-1.0.1.tar.gz
echo ','
$DIRECTORY/build_plugin_entry.sh https://github.com/mattermost/mattermost-plugin-zoom/releases/download/v1.0.7/zoom-1.0.7.tar.gz
echo ','
$DIRECTORY/build_plugin_entry.sh https://github.com/mattermost/mattermost-plugin-jira/releases/download/v2.0.6/mattermost-plugin-jira-v2.0.6.tar.gz
echo ','
$DIRECTORY/build_plugin_entry.sh https://github.com/mattermost/mattermost-plugin-autotranslate/releases/download/v0.1.2/autotranslate-0.1.2.tar.gz
echo ','
$DIRECTORY/build_plugin_entry.sh https://github.com/mattermost/mattermost-plugin-profanity-filter/releases/download/v0.1.0/mattermost-profanity-filter.tar.gz
echo ','
$DIRECTORY/build_plugin_entry.sh https://github.com/mattermost/mattermost-plugin-welcomebot/releases/download/v1.1.0/com.mattermost.welcomebot-1.1.0.tar.gz
echo ','
$DIRECTORY/build_plugin_entry.sh https://github.com/mattermost/mattermost-plugin-jenkins/releases/download/V0.0.3/jenkins-0.0.3.tar.gz
echo ','
$DIRECTORY/build_plugin_entry.sh https://github.com/mattermost/mattermost-plugin-antivirus/releases/download/v0.1.0/antivirus-0.1.0.tar.gz
echo ','
$DIRECTORY/build_plugin_entry.sh https://github.com/mattermost/mattermost-plugin-walltime/releases/download/0.1.1/com.mattermost.walltime-plugin-0.1.1.tar.gz
echo ','
$DIRECTORY/build_plugin_entry.sh https://github.com/mattermost/mattermost-plugin-custom-attributes/releases/download/v1.0.0/com.mattermost.custom-attributes-1.0.0.tar.gz
echo ','
$DIRECTORY/build_plugin_entry.sh https://github.com/mattermost/mattermost-plugin-skype4business/releases/download/v0.1.2/skype4business-0.1.2.tar.gz
echo ','
$DIRECTORY/build_plugin_entry.sh https://github.com/mattermost/mattermost-plugin-aws-SNS/releases/download/v1.0.2/com.mattermost.aws-sns-1.0.2.tar.gz
echo ','
$DIRECTORY/build_plugin_entry.sh https://github.com/mattermost/mattermost-plugin-gitlab/releases/download/v0.3.0/com.github.manland.mattermost-plugin-gitlab-0.3.0.tar.gz
echo ','
$DIRECTORY/build_plugin_entry.sh https://github.com/mattermost/mattermost-plugin-nps/releases/download/v1.0.3/com.mattermost.nps-1.0.3.tar.gz
echo ']'
