### 0.5.0 (TBD)

#### Features
- Add custom date functionality for stashes - [Issue #251](https://github.com/sensu/uchiwa/pull/251)
- Aggregates support - [Issue #255](https://github.com/sensu/uchiwa/pull/255)
- Use JSON Web Tokens (JWT) instead of HTTP basic authentication and implement a login page - [Issue #8](https://github.com/sensu/uchiwa-web/pull/8)

#### Bug Fixes
- Catch possible exception while asserting the check name - [Issue #259](https://github.com/sensu/uchiwa/pull/259)
- Display images when an URL has a port number - [Issue #2](https://github.com/sensu/uchiwa-web/pull/2)
- Hide events with no client and no check - [Issue #9](https://github.com/sensu/uchiwa-web/pull/9)
- Client names may not be unique - [Issue #11](https://github.com/sensu/uchiwa-web/pull/11)

### 0.4.1 (2015-01-16)

#### Features
- Display the stash type on the stashes view (PR #249)

#### Bug Fixes
- Catch any error while asserting event attributes (Issue #236)
- Reuse http.Client on gosensu library (Issue #239)
- Verify the boolean type on richOutput function (Issue #247)

### 0.4.0 (2014-12-08)

#### Features
- Bulk actions support for events & clients (PR #201)
- Display an error page if a client is not found (Issue #200)
- Use `$interval` instead of `pollingFactory` (PR #215)

#### Bug Fixes
- Make sure to always close the connection to Sensu APIs (Issue #211)
- Hide clients name overflow in clients view (Issue #210)
- Fix design of modal window with dark theme (Issue #206)
- Allow Uchiwa to be ran behind a reverse proxy with a location (Issue #197)
- Display missing stash info from client view (PR #218)
- Apply filters to clients & events with bulk selections (PR #214)

### 0.3.4 (2014-11-25)

#### Bug Fixes
- favicon.ico was empty or corrupt
- Imagey regex won't match hostnames or IP addresses

### 0.3.3 (2014-11-21)

#### Features
- Remove imagey-filter so http images do not use Google images proxy

### 0.3.2 (2014-11-04)

#### Features
- Allow for silenced checks to be hidden (Issue #176)

#### Bug Fixes
- Support invalid certificate on API for stash and client deletion

### 0.3.1 (2014-11-03)

#### Bug Fixes
- Display the right check model and rich output when check returns 0
- Support API authentication for stash and client deletion

### 0.3.0 (2014-11-02)

#### Features
- Remove all WebSocket dependency (Issue #127)
- Backend refactoring in Go (Issue #127)
- Improve overall performance and stability
- Add alert badge in navbar when a datacenter is missing
- Stashes dropdown replaced with a modal dialog
- Add custom stash messages, and display them in stashes view (PR #158)
- Add links to DC hover menu (PR #152)
- Display data centers in alphabetic order (PR #153)
- Prettify JSON objects in client & event data (PR #170)

#### Migration Notes
- Backend has been rewritten in Go. `node app.js` or `npm start` commands no longer works.
- Make sure to run `npm install` when installed from source.
- Integer values **timeout** and **refresh** within configuraton file are now handled as seconds instead of milliseconds. Any values >= 1000 will be converted to seconds during runtime.
- Pages */health/[sensu|uchiwa]* now return the associated object content instead of the object itself. */health* is not impacted.

### 0.2.6 (2014-10-08)

#### Features
- Allow client checks to be ordered by history, name, output and time (PR #149)

#### Bug Fixes
- Force toastr position when cookie is missing
- Remove HTML tags from toastr notifications

### 0.2.5 (2014-10-08)

#### Features
- Rich data output for links and images (Issue #86)
- Remove jQuery dependency (PR #132)
- Use angular-toastr module instead of toastr library (PR #132)
- Add real favicon (PR #145)

#### Bug Fixes
- Prevent XSS attacks through toastr library

### 0.2.4 (2014-09-29)

#### Features
- Use source property in events for masquerading (PR #134)
- Date timezone is now determined by the browser (PR #124)
- Move most of sensu.js library logic to AngularJS (PR #124)

#### Bug Fixes
- Avoid HTTP 500 errors on /health page (Issue #128)
- Improve stability when dealing with retrieved data (Issue #119)

### 0.2.3 (2014-09-04)

#### Bug Fixes
- Fix authentication (PR #117)
- Prevent crash when no checks are received (PR #115)

### 0.2.2 (2014-09-02)

#### Features

- Automatic permalinks based on search filters (PR #111)
- Configurable date/time formatting (PR #103)
- Add /health page (PR #108)
- Uchiwa logs are now in JSON format (PR #109)
- Mark active page in sidebar
- Use Socket.IO 1.0 (PR #99)
- Use Express 4 (PR #109)
- Accessibility improvements for status circles (PR #105)
- Fix pill border overflow of datacenters list (PR #97)
- Better unit tests coverage (PR #101)
- Enable Travis CI

#### Bug Fixes

- Fix 'show all' option in clients view
- Display clients with no subscriptions (PR #104)

#### Migration Notes
- Make sure to run `npm install`

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
