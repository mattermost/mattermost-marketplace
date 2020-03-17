# Contributing

Thank you for your interest in contributing! Join the [Plugin Marketplace channel](https://community.mattermost.com/core/channels/plugins-marketplace) on the Mattermost Community Server for discussion about the Plugin Marketplace.


## Reporting issues

If you think you found a bug within the Marketplace code, [please use the GitHub issue tracker](https://github.com/mattermost/mattermost-marketplace/labels/issues/new) on this repository to open an issue. Bugs in the Marketplace in-product experience should be reported on the [Mattermost Server repository](https://github.com/mattermost/mattermost-server/issues/new). Please report bugs within specific plugins on their respective issue trackers.


## Community plugin

To add your plugin to the Marketplace, please open [an issue using this template](https://github.com/mattermost/mattermost-marketplace/issues/new?template=add_plugin.md). To update your plugin, please also open an issue [using this template](https://github.com/mattermost/mattermost-marketplace/issues/new?template=update_plugin.md).


### Playbook

The following is a playbook for use by Mattermost core committers.

#### New plugin

After a new plugin has been submitted, the assignee of the issue posts the following message:
```
## Process checklist
- [ ] Create a private fork under the Mattermost organization; `master` should only contain a `README.md`.
- [ ] Give submitter read access to the fork.
- [ ] Create PR to merge upstream into `master`.
- [ ] Request reviews.
- [ ] [Cut plugin release](https://developers.mattermost.com/internal/plugin-release-process).
- [ ] Add release to Marketplace.
- [ ] Reach out on [Marketing Channel](https://community-release.mattermost.com/private-core/channels/marketing) to tweet about the plugin.
- [ ] Work with `@hanna.park` regarding swag.
```

#### Update plugin

After an update for an existing plugin has been submitted, the assignee of the issue posts the following message:
```
## Process checklist
- [ ] Create PR to merge changes from upstream into `master`.
- [ ] Request reviews.
- [ ] [Cut plugin release](https://developers.mattermost.com/internal/plugin-release-process).
- [ ] Add release to Marketplace.
```
