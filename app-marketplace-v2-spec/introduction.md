# Introduction

  * [Problem Statement](#problem-statement)
  * [Project Summary](#project-summary)
  * [Objectives](#objectives)
    + [Full Backwards Compatibility](#full-backwards-compatibility)
    + [Supported Use Cases](#supported-use-cases)
    + [Launch App Marketplace](#launch-app-marketplace)
    + [Scoped **Cloud App** API access](#scoped---cloud-app---api-access)
    + [Meet **Mattermost Cloud** Operational
      Requirements](#meet---mattermost-cloud---operational-requirements)
  * [Technical Principles](#technical-principles)

## Problem Statement

The current extensibility framework supports http-based **Integrations** and
Go/React **Plugins**.

1. **Integrations** are based on simple HTTP1.1 messages, and provide very basic
   interactive functionality, and are not easily hostable on-prem. "As is" this
   framework is not sufficient to develop most interactive in-place use-cases.
   They are mostly supported on mobile.
2. **Plugins** do not offer significant capabilities on mobile. Historically,
   they have been implemented with WebApp first in mind.
3. **Plugins** require development in Go, React; expensive and slow to develop
   and review.
4. **Plugin** installation and configuration management UX is neither
   consolidated, nor consistent.
5. **Plugins** are not sufficiently isolated neither in terms of data access,
   nor in terms of performance for cloud (and some enterprise) deployments.
6. **Slack** compatibility has not been actively maintained, is behind.

## Project Summary 

1. Build the **App Marketplace**, the place to discover,
   install/upgrade/downgrade, and configure **Apps** from.  **App Marketplace**
   will be largely based on the existing **Plugin Marketplace**, but needs to be
   significantly extended in functionality. Specifically, supporting
   **Integrations** and **Cloud Apps**, as well as in-place configuration - need
   to be implemented.
2. **App Listing** will be a new single, installable/configurable listing,
   combining the elements of **Integrations** (webhooks, commands), **Plugin**,
   and **Cloud App** into a single deployable/configurable package.
3. **Cloud Apps** will be http-based, interactive services capable of
   implementing content sharing, in-place viewing and editing integration with
   other 3rd party services. **Cloud Apps** will be designed as "mobile first".
   In spite of the name, **Cloud Apps** the long-term objective is to support
   **Cloud Apps** in on-prem deployments.
4. **Cloud Apps** will likely benefit from new services/capabilities
   added/refactored in the core Mattermost product. #TODO. 
5. Launch and migration: #TODO.

## Objectives
### Full Backwards Compatibility
1. Existing Integrations and Plugins must continue to work with no changes.
2. The UX exposure of the Plugin Marketplace, Integrations Directory, Plugin
   Settings may be changed.

### Supported Use Cases 
1. **UX Hooks** - Invite channel members to a call, from a UX element in channel
2. **Dynamic interactive data inputs** - Share a post/thread upstream from the
   Post menu, specify additional (dynamic) metadata interactively
3. Interactive **subscription notifications** (a.k.a. incoming webhooks) -
   Subscribe to Google Calendar events, with *Accept/Decline with Note* and
   *Propose a Different time* from the notification message
4. **Dashboard** - a widget (LHS/RHS equiv) #TODO

### Launch App Marketplace
1. Launch with 10+5 New **Cloud Apps**. #TODO list
2. Launch as a full set of **Plugins**, **Integrations**, and **Cloud Apps**.
3. Interactive (cloud-apps only) in-place install and configuration, integrated
   into the marketplace UX.

### Scoped **Cloud App** API access
1. Introduce capability scope declaration and enforcement for **Cloud Apps**
   #TODO

### Meet **Mattermost Cloud** Operational Requirements
1. Hostable as multi-tenant by the upstream application provider, or as
   single-tenant in Mattermost Cloud.
2. In Mattermost Cloud, packaged for fast deployment (k8s service?,
   serverless?).
3. In Mattermost Cloud, 95% 100ms target response time (first-byte sent to last
   byte received) for user actions (in addition to the upstream R/T time).
4. Mattermost Cloud-first design for monitoring, logging, introspection
5. #STRETCH on-prem hostable, maybe K8s only.
6. Storage, other service dependencies as a "standard" set of services shared
   across all apps, but not across MM instances. Encourage the use of props in
   MM. data model?
7. Versioning and CI? #TODO

## Technical Principles
1. All user actions are functions (basicallu, commands)
2. MM->App use JWT, App->MM use OAuth2.
3. Encourage the use of props on User, Channel, Post (any other entities?)
4. "Standard" set of dependency services, deployed per-instance; 