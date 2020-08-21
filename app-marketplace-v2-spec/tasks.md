# Tasks

* [Commands](#commands)
   + [Support JSON in/out](#support-json-in-out)
   + [Dynamic argument definition](#dynamic-argument-definition)
   + [Autocomplete types](#autocomplete-types)
   + [Command box](#command-box)
   + [Embedded commands](#embedded-commands)
* [Interactive Dialog](#interactive-dialog)
* [Bot Conversation](#bot-conversation)
* [Storage](#storage)
* [Mattermost Cloud](#mattermost-cloud)
   + [Packaging and Deployment](#packaging-and-deployment)
   + [Monitoring](#monitoring)
* [Mattermost Server](#mattermost-server)
   + [Plugin API](#plugin-api)

## Commands

### Support JSON in/out
- #TODO design JSON/flag mapping
Command JSON out design:
- Should identify the source / location / type (e.g. POST_ACTION, CHANNEL_HEADER, SETTINGS_MENU). This could be inferred by the request URL, but we add it to allow several types to target the same URL.
- Should identify the context of the command. It is highly dependent on the type of the command.
  - user_id will be present on most cases
  - POST_ACTION requires the post_id. To facilitate the work on the service provider, we may want to add channel_id and team_id.
  - CHANNEL_HEADER requires the channel_id. To facilitate the work on the service provider, we may want to add the team_id.
  - SLASH_COMMAND requires the channel_id and the content of the slash command.
  - #TODO design all types and their contexts
- If we move forward with the client state idea, it should send the client state

### Dynamic argument definition
- Know what options are relevant in the context

### Autocomplete types
- Date picker
- Plain text (many versions?)
- ...

### Command box 
- Input with (filtered) autocomplete, pre-configured initial state
- Display markdown output

### Embedded commands
- Button
- Select
- On/Off Checkbox

## Interactive Dialog
- Mutable upon entering/selecting fields
- Dynamic list support

## Bot Conversation
- Build on the Settings Panel

## Storage
- Research tradeoffs of relying on props on User, Channel, Team
- Still need persistence, cache with invalidation, background jobs, stateful timers (reminders)
- Should Cloud Apps using 3rd paty services explicitly declare them? Enforcement?

- Dani Q: Shouldn't the service provider handle the storage? If the service provider is a mattermost plugin, will use KVStore, if it is an external service, they can handle the storage. Background jobs and timers should be done again in the service provider, and the action they need to perform to do it using the API or incoming webhooks. The only "storage" needed IMO is client state, for which I have already some ideas.

## Mattermost Cloud
### Packaging and Deployment 
### Monitoring

## Mattermost Server

### Plugin API
- Is StorePluginConfig robust enough? Support individual keys? #TODO research, how else would Cloud Apps store config?
