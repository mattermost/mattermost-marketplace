# Architecture.

  * [**Cloud Apps**.](#--cloud-apps--)
    + [Overview.](#overview)
    + [Authentication with Mattermost.](#authentication-with-mattermost)
    + [Authentication with Upstream Applications.](#authentication-with-upstream-applications)
    + [Commands and interactivity.](#commands-and-interactivity)
    + [Event notifications.](#event-notifications)
    + [Installation.](#installation)
    + [Configuration.](#configuration)
    + [mattermost-cloud-app.json](#mattermost-cloud-appjson)
  * [mattermost-plugin-cloud-apps](#mattermost-plugin-cloud-apps)
  * [mattermost-server](#mattermost-server)
    + [Minimal Changes to the Core Code.](#minimal-changes-to-the-core-code)
  * [Cloud App Marketplace](#cloud-app-marketplace)

## **Cloud Apps**.

### Overview.
1. Mattermost interacts with **Cloud Apps** in 2 ways: **Commands** that have request/response semantics, and 1-way **Event notifications**.
  - Dani Q: Alternatively, **Event notifications** can be sent to the clients, and the clients trigger commands if needed.
2. Cloud Apps interact with Mattermost primarily by using the Mattermost REST APIs. The "Response" functionality of command handlers will be greatly simplified for Cloud Apps.
  - Dani Q: Response == Integrations 2.0 document's "Side Effects". Should we unify naming? Do you want me to add in this document the list I made of possible future side effects?

### Authentication with Mattermost.

1. When a **Cloud App** is installed, it receives "bot credentials" that it can
   use to invoke Mattermost REST APIs as a bot.
   - Q: can this leverage bot accounts?
   - Dani Q: Should this be a special bot with scoped functionality? Since cloud app could potentially be user level (I as an user install it so only I see it in my client) or admin level (I as a sysadmin install it so all users on my instance can use it), not sure at what level bot accounts are.
1. When a **Cloud App** is installed, it gets a shared secret it can use to
   decode a Mattermost JWT.
   - Dani Q: Involve security already to give an opinion over this?
1. Commands sent to the **Cloud App** include a JWT. The JWT contains a
   user-scoped Mattermost REST API token. It should then use the token to act on
   behalf of the user when posting back to Mattermost.
   - Q: What if the Cloud App wants to post back as the bot? Seems no problem,
     just use the correct API token.
1. Post-event notifications sent to the **Cloud App** include a JWT. The JWT
   contains a bot-scoped Mattermost REST API token.

### Authentication with Upstream Applications.
1. #TODO.
  - Dani Q: Same as current outgoing webhooks?

### Commands and interactivity.
1. A **Command** is a fundamental unit of executing specific user instructions
   by a **Cloud App**. They are "embeddable" as Post Actions, and "interactive"
   as Interactive Dialogs.
2. The protocol is to be simplified, so **Commands** return a simple response
   message to instruct immediate client action: e.g. open an Interactive Dialog
   or go to a URL. Creating and updating Posts is to be done via the REST API.
3. There are 3 kinds of commands:
- [Slash commands](commands.md#slash-commands) are entered from the message box, with autocomplete.
- [Embedded commands](commands.md#embedded-commands) are what was previously known as Post Actions or Interactive Messages. They can be embedded in Post Slack attachments, and #STRETCH markdown. **Buttons** launch fully-configured commands, **checkboxes** toggle boolean flags, and **selects** allow to choose a value from a list.
  - Dani Q: This will include other UI elements, like what I have been working for mobile, right? (channel header button, channel header menu, main menu ...)
- [Interactive commands](commands.md#interactive-commands) are Interactive Dialogs, or Bot Conversations.

- Dani Q: Should we clarify whether the locations registered are shown to the user during the scope confirmation? In the sense of: "my app register these slash commands: /myapp, /myappanother. It also register these UI elements: 2 post actions (Webapp and mobile), 1 channel header button (webapp), 1 channel menu (mobile), 1 LHS component (webapp). Do you accept?" 

See [Commands and Interactivity](commands.md) for more.

### Event notifications.
- MessagePosted (matches the existing webhook)
  Must be scoped to: User, Channel.
- ChannelCreated
- UserCreated
- UserJoinedChannel
- UserJoinedTeam
- UserLeftChannel
- UserLeftTeam
- UserUpdated

- Dani Q: Should we specify that events are also something "scoped" and that will be informed and accepted by the person installing the app? In the sense of: "my app will register the following events: UserCreated, UserJoinedTeam. Do you accept?". And only those events will be sent to that integration.

### Installation.
- Performs all necessary registrations and generates/stores secrets
- Run the "install" command

### Configuration. 
- Expose pre-configured commands

### Unninstallation.
- Remove all registrations from the app
- Deauthorize all the tokens provided to the app
- Inform the app endpoint that the application has been unninstaled and deauthorized

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

