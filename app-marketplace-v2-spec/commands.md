# Commands and Interactivity.

* [Slash Commands.](#slash-commands)
* [Embedded Commands.](#embedded-commands)
* [Interactive Commands.](#interactive-commands)
* [Command Request.](#command-request)
* [Command Request Expansion:](#command-request-expansion-)
* [Command Response.](#command-response)

All interactive user actions in Cloud Apps result in executing commands. 

- [Slash commands](#slash-commands) are entered from the message box, with
  autocomplete.
- [Embedded commands](#embedded-commands) are what was previously known as Post
  Actions or Interactive Messages. They can be embedded in Post Slack
  attachments, and #STRETCH markdown. **Buttons** launch fully-configured
  commands, **checkboxes** toggle boolean flags, and **selects** allow to choose
  a value from a list.
- [Interactive commands](#interactive-commands) are Interactive Dialogs, or Bot
  Conversations.

## Slash Commands.
**Commands** are traditionally submitted from the message box, by typing a
`/<trigger>`.
- Autocomplete functionality will be extended to improve dynamic definitions and
  new data types (date/time, plain text, json, file, etc).
- A "command box" component will be available to allow executing an isolated
  command.

## Embedded Commands.
- Post Actions represent a way of encoding 1-click, fully pre-configured
  commands into Posts (Slack Attachments).
- The functionality will be adjusted to match slash-command autocomplete
- Can we implement "embedded interactive dialogs" to submit several fields at
  once?
- #TODO spec: slack attachment format

## Interactive Commands.

### Types of Interactive Commands
1. Interactive dialogs.
   1. Can be directly launched from all relevant UX locations.
   1. Can be launched as a result of another command.
   1. Can be pre-configured with initial set of fields/values.
   2. Can dynamically fetch relevant field data and reconfigure, including the
      set of fields displayed, based on the user inputs (**Autocomplete Query
      API**).
   3. In the end, submits a command (or Cancel).

2. Bot Conversations.
   1. Can be directly launched from all relevant UX locations.
   1. Can be launched as a result of another command.
   1. Navigates the user to the DM with the bot, continues the conversation
      there. (can this be in a separate modal??)
   2. Steps are either static or dynamic.
   3. In the end, submits a command (or Cancel).

### Invoking Interactive Commands
- Declared as interactive, can be invoked as slash commands
- Can act as commands associated with Embedded Buttons - clicking on a button
  would launch an interactive command, rather than submit a message to the
  server.
- Can be specified as an action for extensible UX elements: Channel Header
  buttons and menu, Post menu, main menu.

## Command Request.
Legacy HTTP slash commands are sent as HTTP POSTs with form encoding, see [Slash
Commands - Basic
Usage](https://developers.mattermost.com/integrate/slash-commands/#basic-usage).

V2 HTTP commands will be sent as JSON, with an expandable Context. Example:

```http
POST / HTTP/1.1
Accept-Encoding: gzip
Accept: application/json
Authorization: Token nezum4kpu3faiec7r7c5zt6tfy     #TODO use JWT
Content-Length: xxx
Content-Type: application/json
Host: 127.0.0.1:8080
User-Agent: mattermost-5.xxx
```
```json
{
    "command":{
        "namespace":"jira",
        "function":"create_issue",
        "raw":"/jira create issue --project=MM ...",
        "encoded":{
            "project":"MM",
            "...":"..."
        }
    },
    "security":{
        "token":"xxx #TODO use JWT",
        "trigger_id":"xxx #TODO use JWT"
    },
    "context":{
        "source":{
            "id":"post_menu_xxx",
            "client_or_session_id_to_send_websockets_to":"xxx",
            "props":{}
        },
        "config":{
           "settingX":"value"
        },
        "channel":{},
        "post":{},
        "parent_post":{},
        "root_post":{},
        "team":{},
        "acting_user":{},
        "mentioned":{}
    },
}
```
## Command Request Expansion: 
- #TODO.

## Command Response.
- #TODO.
