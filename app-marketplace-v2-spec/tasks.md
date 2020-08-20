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

## Mattermost Cloud
### Packaging and Deployment 
### Monitoring

## Mattermost Server

### Plugin API
- Is StorePluginConfig robust enough? Support individual keys? #TODO research, how else would Cloud Apps store config?


