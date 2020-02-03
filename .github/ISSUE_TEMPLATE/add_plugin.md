---
name: Add plugin
about: Request to add your plugin to the Marketplace
title: "Add $REPOSITORY_NAME to Marketplace"
labels: Plugin/New
assignees: hanzei

---
<!--
Thank you very for submitting your Plugin! It will go through a review process to make sure it follow the quality standard  of the Marketplace. This process might take a couple of weeks bepending on how many changes are needed.
Read https://developers.mattermost.com/extend/plugins/community-plugin-marketplace/ before submitting your plugin.
-->

#### Summary
<!--
A brief description what your plugin does.
-->

#### Review Commit
<!--
Please link the commit or release that should be used for review.
-->

## Checklist
<!--
Please go trough this checklist and confirm every item. If your plugin doesn't fulfil every item, leave a comment explaining why and if you will fix this.
-->

**Product Requirements**

- [ ] The plugin is published under an Apache v2 compatible license (e.g. no GPL, APGL). A list of compatible licenses can be found [here](https://apache.org/legal/resolved.html#category-a).
- [ ] The source code is available in a public git repository.
- [ ] There is a public issue or bug tracker for the plugin, which is linked in the plugin documentation and linked via `support_url` in the manifest.
- [ ] The plugin provides detailed usage documentation with at least one screenshot of the plugin in action, list of features and a development guide. This is typically a README file or a landing page on the web. The link to the documentation is set as `homepage_url` in the manifest. A great example is the [README of the GitHub plugin](https://github.com/mattermost/mattermost-plugin-github/blob/master/README.md).
- [ ] For every release a changelog has to be publish. The link to it is noted via `release_notes_url` in the manifest.
- [ ] The plugin has to be out of Beta and be released with at least v1.0.0.
- [ ] All configuration has to be possible using the UI of Mattermost.
- [ ] The plugin id defined in the manifest must not collide with the id of an existing plugin in the marketplace. It should follow [the naming convention](https://developers.mattermost.com/extend/plugins/manifest-reference/#id).

**Technical Requirements**

- [ ] The plugin works for 60k concurrent connections and in a high availability environment. (There are currently no tools available to verify this property. Hence, it is checked via code review by a developer)
- [ ] The plugin logs important events on appropriate log levels to allow system admins to troubleshoot issues.

**Security Requirements**

- [ ] The plugin does not expose a vulnerability.
- [ ] The plugin does not include favor the author of the plugin or a third party excessively by e.g. including a bitcoin miner that mines on behalf of the author.
- [ ] The plugins author must notify Mattermost about any vulnerabilities in the future.

**Functional Requirements**

- [ ] The plugin works as expected with the latest version of Mattermost.
- [ ] The plugin works as expected with the latest ESR version of Mattermost. This must not be checked if `min_server_version` is higher than the latest ESR version.
