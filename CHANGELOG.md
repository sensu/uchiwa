### 0.2.1 (2014-08-07)

#### Features
- Navbar icon now links to related page

#### Bug Fixes
- Perform a deep clone for public config display (Issue #78)

### 0.2.0 (2014-08-07)

#### Features
- New user interface! (Issue #55)
- Temporarily silence an element
- Filter and order by attributes
- Display a limited number of elements by default, to reduce page size
- Display custom attributes of checks/clients/events (Issue #58)
- Add an overview of each DC in the navbar
- Ability to link to a client and an incident (Issue #59)
- Filter clients by subscriptions and ability to link it
- Add an info page to display Sensu and Uchiwa basic information
- Dynamic page title (Issue #70)
- Optimize dark theme
- Change licence to MIT (same as Sensu)

#### Bug Fixes
- Better handling of unknown elements (Issue #59)
- Display proper information concerning check details (Issue #72)
- Client event might have shown wrong data
- Validate and initialize missing configuration for Sensu endpoints and Uchiwa

### 0.1.7 (2014-07-30)

#### Features
- The configuration file now use a standard JSON file (see migration notes below) (PR #66)
- Add Sass Grunt task (PR #67)

#### Migrating from 0.1.6 to 0.1.7

The configuration file is now a standard *JSON* file and therefore, has been renamed from **config.js** to **config.json**.

If you already have a **config.js** file, you can still force uchiwa to use it with, for example, the following command: `node app.js -c ./config.js` or by modifing the init script if it was installed from the packages.

Refer yourself to the **config.json.example** file in doubt.

### 0.1.6 (2014-07-28)

#### Bug Fixes
- Truncate checkout output in client modal for dark theme (Issue #57)
- Use underscore .each in utilityService
- Show proper clients & events counts in dashboard panels

### 0.1.5 (2014-07-22)

#### Features
- Use Docker build repository

### 0.1.4 (2014-07-22)

#### Features
- Use AngularJS Routing

### 0.1.3 (2014-06-15)

#### Features
- Manage 3rd party libraries with Bower
- Create Sass themes

### 0.1.2 (2014-06-12)

#### Features
- Updated font-awesome; use database icon for stashes.

#### Bug Fixes
- Stashes can have non-silence paths, use stash.path in stashesService

### 0.1.1 (2014-06-12)

#### Features
- CLI argument parsing for config

#### Bug Fixes
- Do not crash when Sensu return an empty object in dc.js

### 0.1.0 (2014-06-06)

#### Features
- Support multiple Sensu APIs
- Delete a client from the client view

#### Bug Fixes
- Improved error logging
- Display event.action instead of event.flapping in client view

### 0.0.4 (2014-05-31)

#### Features
- Support the upcoming release of Sensu 0.13.0

#### Bug Fixes
- Keep check details expanded in client view

### 0.0.3 (2014-05-28)

#### Features
- Use AngularJS in the frontend
- Add graphic of stashes and events
- Filter events, clients, stashes and checks

#### Bug Fixes
- Google font now protocol-agnostic

### 0.0.2 (2014-05-15)

#### Features
- Add documentation for running uchiwa behind a Nginx proxy

#### Bug Fixes
- Add .map files for javascript librairies
- Better handling of Sensu API path
- Add configuration value for HTTPS
- Use relative path for ressources and links
