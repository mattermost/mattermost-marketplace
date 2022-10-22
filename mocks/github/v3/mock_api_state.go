package v3

import "time"

/* Mock for the MatterMost API state */
var mockMMApiState mockMMApiStateType = mockMMApiStateType{
	"mattermost": map[string][]releaseDetails{
		"mattermost-plugin-antivirus": {
			{"v0.1.0", "v0.1.0", true, false, []releaseAssetDetails{
				{"antivirus-0.1.0.tar.gz", time.Date(2018, 10, 16, 13, 57, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2018, 10, 16, 13, 56, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2018, 10, 16, 13, 56, 0, 0, time.Local)}}},
			{"v0.1.1", "v0.1.1", false, false, []releaseAssetDetails{
				{"antivirus-0.1.1.tar.gz", time.Date(2019, 8, 8, 19, 17, 0, 0, time.Local)},
				{"antivirus-0.1.1.tar.gz.sig", time.Date(2019, 11, 28, 12, 39, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2019, 8, 8, 19, 13, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2019, 8, 8, 19, 13, 0, 0, time.Local)}}},
			{"v0.1.2", "v0.1.2", false, false, []releaseAssetDetails{
				{"antivirus-0.1.2.tar.gz", time.Date(2020, 1, 12, 21, 0, 0, 0, time.Local)},
				{"antivirus-0.1.2.tar.gz.sig", time.Date(2020, 1, 15, 11, 59, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2020, 1, 12, 20, 58, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2020, 1, 12, 20, 58, 0, 0, time.Local)}}}},
		"mattermost-plugin-autolink": {
			{"v0.1.0", "v0.1.0", true, false, []releaseAssetDetails{
				{"mattermost-plugin-autolink-darwin-amd64.tar.gz", time.Date(2018, 6, 4, 20, 9, 0, 0, time.Local)},
				{"mattermost-plugin-autolink-linux-amd64.tar.gz", time.Date(2018, 6, 4, 20, 9, 0, 0, time.Local)},
				{"mattermost-plugin-autolink-windows-amd64.tar.gz", time.Date(2018, 6, 4, 20, 9, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2018, 6, 4, 20, 8, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2018, 6, 4, 20, 8, 0, 0, time.Local)}}},
			{"v0.5.0", "v0.5.0", true, false, []releaseAssetDetails{
				{"mattermost-autolink-0.5.0.tar.gz", time.Date(2019, 4, 11, 20, 14, 2, 0, time.Local)},
				{"Source code.zip", time.Date(2019, 4, 11, 20, 13, 58, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2019, 4, 11, 20, 13, 58, 0, time.Local)}}},
			{"v1.0.0", "v1.0.0", false, false, []releaseAssetDetails{
				{"mattermost-plugin-autolink-1.0.0.tar.gz", time.Date(2019, 6, 11, 10, 4, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2019, 6, 3, 11, 42, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2019, 6, 3, 11, 42, 0, 0, time.Local)}}}},
		"mattermost-plugin-aws-SNS": {
			{"v0.1.0", "v0.1.0", true, false, []releaseAssetDetails{
				{"com.cpanato.aws-sns-0.1.0.tar.gz", time.Date(2019, 5, 28, 7, 7, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2019, 5, 28, 7, 5, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2019, 5, 28, 7, 5, 0, 0, time.Local)}}},
			{"v1.0.0", "v1.0.0", false, false, []releaseAssetDetails{
				{"com.mattermost.aws-sns-1.0.0.tar.gz", time.Date(2019, 6, 11, 15, 41, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2019, 6, 11, 15, 38, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2019, 6, 11, 15, 38, 0, 0, time.Local)}}},
			{"v1.2.0", "v1.2.0", false, false, []releaseAssetDetails{
				{"com.mattermost.aws-sns-1.2.0.tar.gz", time.Date(2021, 1, 15, 6, 14, 0, 0, time.Local)},
				{"com.mattermost.aws-sns-1.2.0.tar.gz.sig", time.Date(2021, 1, 15, 6, 14, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2021, 1, 15, 6, 12, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2021, 1, 15, 6, 12, 0, 0, time.Local)}}}},
		"mattermost-plugin-custom-attributes": {
			{"v0.0.1", "v0.0.1", true, false, []releaseAssetDetails{
				{"com.mattermost.custom-attributes-0.0.1.tar.gz", time.Date(2019, 3, 22, 10, 29, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2019, 3, 1, 9, 7, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2019, 3, 1, 9, 7, 0, 0, time.Local)}}},
			{"v1.0.0", "v1.0.0", false, false, []releaseAssetDetails{
				{"com.mattermost.custom-attributes-1.0.0.tar.gz", time.Date(2019, 6, 3, 9, 56, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2019, 6, 3, 9, 54, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2019, 6, 3, 9, 54, 0, 0, time.Local)}}},
			{"v1.1.0", "v1.1.0", false, false, []releaseAssetDetails{
				{"com.mattermost.custom-attributes-1.1.0.tar.gz", time.Date(2019, 12, 19, 13, 54, 0, 0, time.Local)},
				{"com.mattermost.custom-attributes-1.1.0.tar.gz.sig", time.Date(2020, 1, 9, 8, 18, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2019, 12, 19, 13, 52, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2019, 12, 19, 13, 52, 0, 0, time.Local)}}}},
		"mattermost-plugin-github": {
			{"v0.0.1 - Initial Release", "v0.0.1", true, false, []releaseAssetDetails{
				{"mattermost-github-plugin-darwin-amd64.tar.gz", time.Date(2019, 8, 9, 7, 54, 0, 0, time.Local)},
				{"mattermost-github-plugin-linux-amd64.tar.gz", time.Date(2019, 8, 9, 7, 54, 0, 0, time.Local)},
				{"mattermost-github-plugin-windows-amd64.tar.gz", time.Date(2019, 8, 9, 7, 54, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2019, 8, 9, 7, 14, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2019, 8, 9, 7, 14, 0, 0, time.Local)}}},
			{"v1.0.0", "v1.0.0", false, false, []releaseAssetDetails{
				{"github-1.0.0.tar.gz", time.Date(2020, 5, 28, 14, 54, 0, 0, time.Local)},
				{"github-1.0.0.tar.gz.sig", time.Date(2020, 5, 28, 14, 55, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2020, 5, 28, 14, 16, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2020, 5, 28, 14, 16, 0, 0, time.Local)}}},
			{"v2.0.0", "v2.0.0", false, false, []releaseAssetDetails{
				{"github-2.0.0.tar.gz", time.Date(2020, 10, 15, 5, 29, 0, 0, time.Local)},
				{"github-2.0.0.tar.gz.sig", time.Date(2020, 10, 15, 5, 30, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2020, 10, 15, 5, 25, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2020, 10, 15, 5, 25, 0, 0, time.Local)}}}},
		"mattermost-plugin-gitlab": {
			{"Mimic github plugin", "0.1.0", false, false, []releaseAssetDetails{
				{"com.github.manland.mattermost-plugin-gitlab-0.1.0.tar.gz", time.Date(2019, 4, 17, 15, 49, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2019, 4, 17, 15, 44, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2019, 4, 17, 15, 44, 0, 0, time.Local)}}},
			{"finishes and more", "v0.2.0", false, false, []releaseAssetDetails{
				{"com.github.manland.mattermost-plugin-gitlab-0.2.0.tar.gz", time.Date(2019, 5, 6, 9, 12, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2019, 5, 6, 9, 3, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2019, 5, 6, 9, 3, 0, 0, time.Local)}}},
			{"Polishing", "v0.3.0", false, false, []releaseAssetDetails{
				{"com.github.manland.mattermost-plugin-gitlab-0.3.0.tar.gz", time.Date(2019, 6, 5, 16, 45, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2019, 6, 5, 16, 42, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2019, 6, 5, 16, 42, 0, 0, time.Local)}}},
			{"v1.0.0", "v1.0.0", false, false, []releaseAssetDetails{
				{"com.github.manland.mattermost-plugin-gitlab-1.0.0.tar.gz", time.Date(2019, 8, 14, 17, 16, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2019, 8, 14, 15, 56, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2019, 8, 14, 15, 56, 0, 0, time.Local)}}}},
		"mattermost-plugin-jenkins": {
			{"v0.0.1", "v0.0.1", true, false, []releaseAssetDetails{
				{"jenkins-0.0.1.tar.gz", time.Date(2019, 4, 16, 11, 42, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2019, 4, 16, 4, 51, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2019, 4, 16, 4, 51, 0, 0, time.Local)}}},
			{"v1.0.0", "v1.0.0", false, false, []releaseAssetDetails{
				{"jenkins-1.0.0.tar.gz", time.Date(2020, 8, 14, 17, 13, 0, 0, time.Local)},
				{"jenkins-1.0.0.tar.gz.sig", time.Date(2020, 11, 28, 12, 38, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2020, 8, 14, 15, 56, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2020, 8, 14, 15, 56, 0, 0, time.Local)}}},
			{"v1.1.0", "v1.1.0", false, false, []releaseAssetDetails{
				{"jenkins-1.1.0.tar.gz", time.Date(2020, 6, 19, 3, 42, 0, 0, time.Local)},
				{"jenkins-1.1.0.tar.gz.sig", time.Date(2020, 6, 19, 3, 42, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2020, 6, 19, 3, 39, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2020, 6, 19, 3, 39, 0, 0, time.Local)}}}},
		"mattermost-plugin-jira": {
			{"0.1: First release as an independent plugin", "v0.1", true, false, []releaseAssetDetails{
				{"mattermost-jira-plugin-darwin-amd64.tar.gz", time.Date(2017, 11, 29, 18, 49, 0, 0, time.Local)},
				{"mattermost-jira-plugin-linux-amd64.tar.gz", time.Date(2017, 11, 29, 18, 49, 0, 0, time.Local)},
				{"mattermost-jira-plugin-windows-amd64.tar.gz", time.Date(2017, 11, 29, 18, 49, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2017, 11, 29, 18, 46, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2017, 11, 29, 18, 46, 0, 0, time.Local)}}},
			{"0.1.1", "v0.1.1", true, false, []releaseAssetDetails{
				{"mattermost-jira-plugin-darwin-amd64.tar.gz", time.Date(2017, 12, 4, 5, 35, 0, 0, time.Local)},
				{"mattermost-jira-plugin-linux-amd64.tar.gz", time.Date(2017, 12, 4, 5, 35, 0, 0, time.Local)},
				{"mattermost-jira-plugin-windows-amd64.tar.gz", time.Date(2017, 12, 4, 5, 35, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2017, 12, 4, 5, 32, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2017, 12, 4, 5, 32, 0, 0, time.Local)}}},
			{"0.1.2", "v0.1.2", true, false, []releaseAssetDetails{
				{"mattermost-jira-plugin-darwin-amd64.tar.gz", time.Date(2017, 12, 5, 18, 41, 0, 0, time.Local)},
				{"mattermost-jira-plugin-linux-amd64.tar.gz", time.Date(2017, 12, 5, 18, 41, 0, 0, time.Local)},
				{"mattermost-jira-plugin-windows-amd64.tar.gz", time.Date(2017, 12, 5, 18, 41, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2017, 12, 5, 18, 38, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2017, 12, 5, 18, 38, 0, 0, time.Local)}}},
			{"1.0.0", "v1.0.0", true, false, []releaseAssetDetails{
				{"mattermost-jira-plugin-darwin-amd64.tar.gz", time.Date(2018, 7, 15, 21, 17, 0, 0, time.Local)},
				{"mattermost-jira-plugin-linux-amd64.tar.gz", time.Date(2018, 7, 15, 21, 17, 0, 0, time.Local)},
				{"mattermost-jira-plugin-windows-amd64.tar.gz", time.Date(2018, 7, 15, 21, 17, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2018, 7, 15, 20, 57, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2018, 7, 15, 20, 57, 0, 0, time.Local)}}},
			{"v1.0.1", "v1.0.1", true, false, []releaseAssetDetails{
				{"mattermost-jira-plugin-darwin-amd64.tar.gz", time.Date(2018, 7, 24, 9, 47, 0, 0, time.Local)},
				{"mattermost-jira-plugin-linux-amd64.tar.gz", time.Date(2018, 7, 24, 9, 47, 0, 0, time.Local)},
				{"mattermost-jira-plugin-windows-amd64.tar.gz", time.Date(2018, 7, 24, 9, 47, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2018, 7, 24, 9, 48, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2018, 7, 24, 9, 48, 0, 0, time.Local)}}},
			{"v1.0.2", "v1.0.2", true, false, []releaseAssetDetails{
				{"com.mattermost.jira-1.0.2.tar.gz", time.Date(2018, 12, 13, 4, 52, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2018, 11, 7, 5, 8, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2018, 11, 7, 5, 8, 0, 0, time.Local)}}},
			{"v1.0.3", "v1.0.3", false, false, []releaseAssetDetails{
				{"jira-1.0.3.tar.gz", time.Date(2018, 12, 18, 2, 37, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2018, 12, 17, 3, 4, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2018, 12, 17, 3, 4, 0, 0, time.Local)}}}},
		"mattermost-plugin-nps": {
			{"v0.0.2-rc4", "v0.0.2-rc4", true, false, []releaseAssetDetails{
				{"com.mattermost.nps-0.0.2.tar.gz", time.Date(2019, 6, 3, 15, 1, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2019, 6, 3, 14, 57, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2019, 6, 3, 14, 57, 0, 0, time.Local)}}},
			{"v1.0.0", "v1.0.0", false, false, []releaseAssetDetails{
				{"com.mattermost.nps-1.0.0.tar.gz", time.Date(2019, 6, 12, 15, 4, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2019, 6, 7, 9, 32, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2019, 6, 7, 9, 32, 0, 0, time.Local)}}},
			{"v1.1.0", "v1.1.0", false, false, []releaseAssetDetails{
				{"com.mattermost.nps-1.1.0.tar.gz", time.Date(2020, 9, 18, 15, 11, 0, 0, time.Local)},
				{"com.mattermost.nps-1.1.0.tar.gz.sig", time.Date(2020, 9, 18, 15, 11, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2020, 9, 17, 18, 1, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2020, 9, 17, 18, 1, 0, 0, time.Local)}}}},
		"mattermost-plugin-webex": {
			{"v1.0.0", "v1.0.0", false, false, []releaseAssetDetails{
				{"com.mattermost.webex-1.0.0.tar.gz", time.Date(2019, 10, 4, 11, 14, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2019, 10, 4, 11, 11, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2019, 10, 4, 11, 11, 0, 0, time.Local)}}},
			{"v1.1.0", "v1.1.0", false, false, []releaseAssetDetails{
				{"com.mattermost.webex-1.1.0.tar.gz", time.Date(2020, 12, 5, 1, 59, 0, 0, time.Local)},
				{"com.mattermost.webex-1.1.0.tar.gz.sig", time.Date(2020, 12, 7, 5, 14, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2020, 12, 5, 1, 57, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2020, 12, 5, 1, 57, 0, 0, time.Local)}}},
			{"v1.2.0", "v1.2.0", false, false, []releaseAssetDetails{
				{"com.mattermost.webex-1.2.0.tar.gz", time.Date(2021, 7, 27, 9, 33, 0, 0, time.Local)},
				{"com.mattermost.webex-1.2.0.tar.gz.sig", time.Date(2021, 7, 27, 9, 33, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2021, 7, 27, 9, 30, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2021, 7, 27, 9, 30, 0, 0, time.Local)}}}},
		"mattermost-plugin-welcomebot": {
			{"v0.1.0", "v0.1.0", true, false, []releaseAssetDetails{
				{"com.mattermost.welcomebot.tar.gz", time.Date(2018, 9, 14, 23, 30, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2018, 9, 14, 23, 29, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2018, 9, 14, 23, 29, 0, 0, time.Local)}}},
			{"v1.0.0", "v1.0.0", false, false, []releaseAssetDetails{
				{"com.mattermost.welcomebot-1.0.0.tar.gz", time.Date(2019, 6, 7, 9, 22, 0, 0, time.Local)},
				{"com.mattermost.welcomebot-1.0.0.tar.gz.sig", time.Date(2019, 11, 21, 10, 36, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2019, 6, 3, 14, 11, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2019, 6, 3, 14, 11, 0, 0, time.Local)}}},
			{"v1.1.0", "v1.1.0", false, false, []releaseAssetDetails{
				{"com.mattermost.welcomebot-1.1.0.tar.gz", time.Date(2019, 6, 17, 13, 7, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2019, 6, 17, 13, 6, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2019, 6, 17, 13, 6, 0, 0, time.Local)}}}},
		"mattermost-plugin-zoom": {
			{"v0.1.0", "v0.1.0", true, false, []releaseAssetDetails{
				{"mattermost-zoom-plugin-darwin-amd64.tar.gz", time.Date(2017, 12, 4, 19, 41, 0, 0, time.Local)},
				{"mattermost-zoom-plugin-linux-amd64.tar.gz", time.Date(2017, 12, 4, 19, 41, 0, 0, time.Local)},
				{"mattermost-zoom-plugin-windows-amd64.tar.gz", time.Date(2017, 12, 4, 19, 41, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2017, 12, 4, 19, 33, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2017, 12, 4, 19, 33, 0, 0, time.Local)}}},
			{"v1.0.0", "v1.0.0", true, false, []releaseAssetDetails{
				{"mattermost-zoom-plugin-darwin-amd64.tar.gz", time.Date(2018, 7, 15, 21, 18, 0, 0, time.Local)},
				{"mattermost-zoom-plugin-linux-amd64.tar.gz", time.Date(2018, 7, 15, 21, 18, 0, 0, time.Local)},
				{"mattermost-zoom-plugin-windows-amd64.tar.gz", time.Date(2018, 7, 15, 21, 18, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2018, 7, 15, 21, 4, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2018, 7, 15, 21, 4, 0, 0, time.Local)}}},
			{"v1.5.0-cloud", "v1.5.0-cloud", false, false, []releaseAssetDetails{
				{"zoom-1.5.0-cloud.tar.gz", time.Date(2020, 11, 4, 12, 7, 0, 0, time.Local)},
				{"zoom-1.5.0-cloud.tar.gz.sig", time.Date(2020, 11, 4, 12, 8, 0, 0, time.Local)},
				{"Source code.zip", time.Date(2020, 11, 4, 12, 4, 0, 0, time.Local)},
				{"Source code.tar.gz", time.Date(2020, 11, 4, 12, 4, 0, 0, time.Local)}}}}}}
