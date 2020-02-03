---
name: Update plugin
about: Request to update your plugin in the Marketplace
title: "Update $REPOSITORY_NAME to $VERSION"
labels: Plugin/Update
assignees: hanzei

---

<!--
Thank you very for continuing to develop and maintain your plugin. It will go through a review process to make sure all requirements are still met since the last release.
-->

#### Summary
<!--
Are there any notable changes since the last release?
-->

#### Review Commit
<!--
Please link the commit or release that should be used for review.
-->

## Checklist
- [ ] [All requirements](https://developers.mattermost.com/extend/plugins/community-plugin-marketplace/#requirements-for-adding-community-plugin-to-the-marketplace) are still met.
- [ ] The release also has to follow semver. This specifically means for plugins:
  - If the plugin exposes an API for inter plugin communication, breaking changes to the API require a major version bump.
  - If an update requires manual migration actions from the administrator, a major version bump is required.
- [ ] A changelog has been published. The link to it is noted via `release_notes_url` in the manifest.
