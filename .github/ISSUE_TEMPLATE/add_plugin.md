---
name: Add plugin
about: Request to add your plugin to the Marketplace
title: "Add $REPOSITORY_NAME to Marketplace"
labels: Plugin/New
assignees: hanzei
---

<!--
Thank you very for submitting your plugin for consideration! A review process is required to ensure your plugin adheres to the quality standard of the Marketplace. This process may take a couple of weeks depending on Mattermost staff availability and any changes that are required.
Read https://developers.mattermost.com/extend/plugins/community-plugin-marketplace/ before submitting your plugin.
-->

#### Summary
<!--
A brief description what your plugin does. Consider including screenshots to help illustrate.
-->

#### Review commit
<!--
Please link to an open source repository and release that should be used for review. As Mattermost code reviews and builds all plugins itself when listing in the Marketplace, the link cannot point at an already-built plugin.
-->

## Checklist
<!--
Go through this checklist and confirm every item. 

It's fine if your plugin doesn't fulfill every item (e.g. it isn't production ready yet). You can still submit it! You can also still do code changes while the plugin is in review and fix issues on the fly.

Even if your plugin isn't production-ready, it might still be added to the Marketplace as "Beta". See https://developers.mattermost.com/extend/plugins/community-plugin-marketplace/#beta-plugins for more details.

If your plugin isn't production-ready, please leave a comment stating whether you plan to fulfill the whole checklist or submit a "Beta" plugin.
-->

**Product requirements**

- [ ] The plugin is published under an [Open Source license](https://opensource.org/licenses/alphabetical).
- [ ] The source code is available in a public Git repository.
- [ ] There is a public issue or bug tracker for the plugin, which is linked in the plugin documentation and linked via `support_url` in the manifest.
- [ ] For the current release and upcoming ones a changelog has to be published, with a link recorded in the `release_notes_url` property of the `plugin.json` manifest.
- [ ] The plugin has to be out of Beta and be released with at least v1.0.0.
- [ ] All configuration is accessible via the Mattermost interface.
- [ ] The plugin ID defined in the manifest must not collide with the ID of an existing plugin in the Marketplace. It should follow [the documentation's suggested naming convention](https://developers.mattermost.com/extend/plugins/manifest-reference/#id).

**Technical requirements**

- [ ] The plugin works for 60k concurrent connections and in a High Availability deployment. **Note:** There are currently no publicly-available tools to verify these properties. As such, they are checked during code review by a developer.
- [ ] The plugin logs important events on appropriate log levels to allow System Admins to troubleshoot issues.

**Security requirements**

- [ ] Security reviews do not reveal any exploitable vulnerabilities in the plugin.
- [ ] The plugin provides an email address or a username on the [Community Server](https://community.mattermost.com) used to report vulnerabilities in the future. Please post it into this issue or send it to ben.schumacher@mattermost.com.

**Functional requirements**

- [ ] The plugin must set a `min_server_version` in the manifest.
- [ ] The plugin must work on all Mattermost versions greater than or equal to the `min_server_version`.

**Documentation requirements (checked by a Technical Writer)**

- [ ] The plugin must include detailed usage documentation with at least one screenshot of the plugin in action, list of features, and a development guide. This is typically a `README` file or a landing page on the web. The link to the documentation is set as `homepage_url` in the manifest. A great example is the [`README` of the GitHub plugin](https://github.com/mattermost/mattermost-plugin-github/blob/master/README.md). Typical components of documentation include:

    * Requirements/Prerequisites
    * Installation steps
    * Configuration steps
    * Usage
    * Troubleshooting
    * Screenshots (if available)
    * Link or email address for help/support

The `plugin.json` file should include formatting consistent with the System Console interface. You can join the [Documentation channel](https://community.mattermost.com/core/channels/documentation) for assistance.
