# Plugin Marketplace

The Plugin Marketplace is a collection of plugins for use with [Mattermost](https://github.com/mattermost/mattermost-server). This repository houses the stateless HTTP service that is run at https://api.integrations.mattermost.com. It is meant to be queried by the Mattermost server to enable plugin discovery by System Admins.

Although Mattermost hosts the Marketplace as an AWS Lambda function backed by S3 and CloudFront, the core feature set is designed for use in any hosting environment, enabling private, self-hosted collections of plugins.

Read more about the [Plugin Marketplace Architecture](https://docs.google.com/document/d/1tVj0eNwMdIIGn8YoTs-cYz9NYvXjqx6bqWH-wa-yDLk/edit).

## Other Resources

This repository houses the open-source components of the Plugin Marketplace. Other resources are linked below:

- [Mattermost the server and user interface](https://github.com/mattermost/mattermost-server)

## Get Involved

- [Join the discussion on ~Plugin Marketplace](https://community.mattermost.com/core/channels/plugins-marketplace)

## Developing

### Environment Setup

1. Install [Go](https://golang.org/doc/install)

### Running

Simply run the following:

```
$ make run-server
```

### Testing

Running all tests:

```
$ make test
```

### Updating plugins.json

To fetch all new release from GitHub, run

```
make plugins.json
```

Make sure to double check the `diff` of `plugins.json` to ensure every release get added correctly.

If you run this command multiple times, GitHub might rate limit you. Using an API token does help:

`export GITHUB_TOKEN=<github token>`

### Deploying as a Lambda Function

In addition to running as a standalone server, the Marketplace is also designed to run as a Lambda function, compiling the `plugins.json` database into the binary for immediate access without further configuration.

### Automatic Deployment

Changes merged to `master` are automatically deployed to https://api.staging.integrations.mattermost.com.

Changes merged to `production` are automatically deployed to https://api.integrations.mattermost.com.

When adding or updating the plugins database (or corresponding tooling), submit the changes directly to `production`, and then merge `production` immediately back to `master` to reduce unnecessary merge conflicts. All other changes should be committed directly to `master`. Changes pending release from `master` should be merged to `production` after qualification and in coordination with any supporting Mattermost server release.

### Manual Deployment

Simply run the following:

```
$ SLS_STAGE=staging make deploy-lambda
```

To iterate quickly after the Cloud Formation stack is up, simply run:

```
$ SLS_STAGE=staging make deploy-lambda-fast
```
