# Architecture.

  * [Cloud Apps.](#cloud-apps)
    + [Authentication with Mattermost.](#authentication-with-mattermost)
    + [Authentication with Upstream Applications.](#authentication-with-upstream-applications)
    + [Cloud App to Mattermost.](#cloud-app-to-mattermost)
    + [Mattermost to Cloud App.](#mattermost-to-cloud-app)
    + [Functions](#functions)
    + [Interactive Inputs.](#interactive-inputs)
      - [Slash Commands](#slash-commands)
      - [Widgets](#widgets)
      - [Dialogs](#dialogs)
    + [Change notifications.](#change-notifications)
    + [Autocomplete lookups](#autocomplete-lookups)
    + [Lifecycle - Install](#lifecycle---install)
    + [Lifecycle Uninstall.](#lifecycle-uninstall)
    + [Configuration.](#configuration)
    + [mattermost-cloud-app.json](#mattermost-cloud-appjson)
  * [mattermost-plugin-cloud-apps](#mattermost-plugin-cloud-apps)
  * [mattermost-server](#mattermost-server)
    + [Minimal Changes to the Core Code.](#minimal-changes-to-the-core-code)
  * [Cloud App Marketplace](#cloud-app-marketplace)

## Cloud Apps.

### Authentication with Mattermost.
1. OAuth2 will be used as a primary mechanism of authentication for the Mattermost REST API.
2. Installing a **Cloud App** onto a Mattermost instance will provision sufficient credential to the app to invoke the REST API as the bot (priviledged? scoped? #TODO), and to verify/decode the JWT on the incoming requests from the Mattermost Server
3. We should enable 3 use-cases:
  a. Cloud App** using a "bot"-level token (to use APIs as the bot)
  b. **Cloud App** using user-level OAuth2 tokens. Will require users' authoriation when first using the app.
  c. "A trusted app" - all users will automatically get "on-behalf" tokens, like with Atlassian Connect. @crspeller mentioned a "trusted app" concept in OAuth2, need to research.
4. Function requests, and Event notifications sent to the Cloud App will include a JWT. The JWT will contain the acting Mattermost UserID (or Bot UserID for non-user-specific notifications)

### Authentication with Upstream Applications.
1. Authentication with upstream applications will be left entirely to the Cloud App
OR
2. Mattermost will provide a "service" for mapping user accounts using OAuth2

### Cloud App to Mattermost.
Cloud Apps interact with Mattermost primarily by using the **Mattermost REST
APIs**. The "Response" functionality of legacy command and action handlers will
be greatly simplified for Cloud App functions. The cloud app functions will
respond as either:
  - 200 OK + (optional) simple markdown + (optional) JSON.
  - Error code + error message.
  - A (302?) "redirect" to navigate to a channel, external browser, an iframe
    popup, or run another function.
  - #TODO @larkox to reconcile with the "side effects" of his Integrations 2.0
    spec.
  - #TODO @levb, @larkox to propose redirect and execute response format.
  - #TODO @crspeller, @lieut-data anything to borrow/learn from the workflow
    projects.

### Mattermost to Cloud App.
Mattermost interacts with Cloud Apps in 3 ways: 
  - **Functions** that are invoked via Interactive Inputs, and have
    request/response semantics. They may be synchronous, or execute in the
    background once invoked.
  - 1-way **Change notifications** that are sent after changes occur (from the
    server).
  - **Autocomplete** queries.

### Functions 
1. Functions are fundamental units of executing specific user instructions
by Cloud Apps. 
2. They are called by Slash Commands, Widgets, Modal and Bot Dialogs. Commands, Modal and Bot dialogs are also callable.
3. Functions may accept both structured JSON inputs, as well as "raw" command strings with options.
4. Dialogs and commands that are fully valid against SlashCommands definition are expected to submit structured data.
5. Widgets "hardcode" their invocation of functions, so they can do as they please.
6. Slash Commands invoke the closest-matched function. If the command string passes validation, it is submitted as JSON, otherwise as raw text. Raw text is always submitted.
7. Expandable context provides functions with most data dependencies pre-fetched, as instructed by the App Manifest, or the function binding.

### Interactive Inputs.
See [Interactive Inputs](commands.md) for more.
  - **Slash Commands** with autocomplete, execute functions.
  - **Widgets** that can be bound to directly execute functions (Post menu items, channel menu, interactive post buttons, etc). In the future they may be extended to become in-place forms or autocomplete boxes.
  - **Modal Dialogs** and **Bot Dialogs**, that gather/edit parameter values, then execute a function.

#### Slash Commands
1. Slash commands will have interactive auto-complete, available in the usual message box in Mattermost. All declared commands that are not "hidden" will be available for autocomplete.
2. #TODO we need to be able to re-define what is visible (valid?) based on the state already entered, and hide the irrelevant options. @iomodo can you take a look at this?
3. #STRETCH a "slash command" widget that would execute a pre-configured autocomplete command, in the given context (team/channel or Bot DM), and display the result.
4. Slash commands may have a --modal option that would launch a modal to gather the inputs, rather than autocomplete. The parameters in the slash command are used as the default values in the modal.
5. All valid slash commands must be declared in the App manifest

#### Widgets
1. A widget is a UX element "hard-wired" to execute a function upon a simple user interaction.
2. Widgets must be pre-declared in the manifest, or inserted in Posts.
3. Widgets' visibility in Mattermost UX can be toggled contextually (user, channel, team prefs) #TODO how?
4. Widgets' visibility in Mattermost UX can be toggled dynamically, via a lookup API #TODO how?
5. Widgets can be stateful (e.g. on/off button)

#### Dialogs
- #TODO

### Change notifications.
1. Subjects
  - *App Lifecycle* Installed/Uninstalled
  - MessagePosted (matches the existing webhook)
      Scoped to: User, Channel.
  - ChannelCreated
  - UserCreated
  - UserJoinedChannel
  - UserJoinedTeam
  - UserLeftChannel
  - UserLeftTeam
  - UserUpdated
2. Expandable context provides functions with most data dependencies pre-fetched, as instructed by the App Manifest, or the function binding.
3. Overall scope in the App Manifest. #STRETCH individual scoped subscriptions for performance - channel, user, etc.

### Autocomplete lookups
- #TODO

### Lifecycle - Install
- Installation is initiated in the App Marketplace, from the Mattermost server. An application secret is generated by the App and must be provided to the Mattermost server.
- Mattermost server completes the handshake with the app by sending it "App Installed" message, which includes the App secret. 
- Once the handshake is complete, the `install` function is called on the App, as configured in the App Manifest

### Lifecycle Uninstall.
- Remove all registrations from the app
- Deauthorize all the tokens provided to the app, invalidate the master secret
- Inform the app endpoint that the application has been unninstaled and deauthorized

### Configuration. 
- A list of widgets

### mattermost-cloud-app.json
- Supports listing and installing/upgrading "standalone" plugins, from the Plugin Marketplace.
- Supports installing "pure" **Cloud Apps**
- Supports installing "pre-declared" legacy webhooks and slash commands (#TODO)

## mattermost-plugin-cloud-apps
- Single-tenant vs Multi-tenant #TODO Discussion
- Install APIs
   + Cloud App - Hooks, etc.
   + 

## mattermost-server

### Minimal Changes to the Core Code.
1. Implement as a "Cloud Apps" plugin (mandatory, pre-loaded).
2. Preserve the current plugin and integrations support "as is"

## Cloud App Marketplace

