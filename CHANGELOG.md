### 1.4.0 (2018-12-14)
#### Features
- Add support for client invalidation - [PR #24](https://github.com/sensu/uchiwa-web/pull/24)

#### Bug Fixes
- Update `mapstructure` dependency to prevent race conditions and segfaults - [Issue #799](https://github.com/sensu/uchiwa/issues/798)

#### Other
- Upgrade dev dependencies

### 1.3.1 (2018-06-11)
#### Bug Fixes
- Fix segfault on unresponsive APIs - [PR #787](https://github.com/sensu/uchiwa/pull/787)
- Fix only one silenced entry appearing - [PR #788](https://github.com/sensu/uchiwa/pull/788)

### 1.3.0 (2018-06-06)
#### Features
- The fetching and building of data from multiple Datacenters will now happen
  concurrently - [PR #781](https://github.com/sensu/uchiwa/pull/781)
  [PR #784](https://github.com/sensu/uchiwa/pull/784)

### 1.2.0 (2018-04-27)
#### Features
- Add confirmation prompt when deleting silenced entry that could affect
  mulitple clients - [Issue #767](https://github.com/sensu/uchiwa/issues/767)
- Silenced entries can now be scheduled (requires Sensu server 1.2.0) - [Issue
  #769](https://github.com/sensu/uchiwa/issues/769)
- Use `PORT` environment variable if available - [Issue
  #772](https://github.com/sensu/uchiwa/issues/772)
- Both `expire_on_resolve` and `expire` can be configured when creating a
  silenced entry - [Issue #776](https://github.com/sensu/uchiwa/issues/776)

#### Bug Fixes
- Properly handle chunked transfer encoding from Sensu API -
[Issue #751](https://github.com/sensu/uchiwa/issues/751)
- Fix select element background color with uchiwa-dark theme - 
[Issue #763](https://github.com/sensu/uchiwa/issues/763)
- Fix dropdown menus height - [PR
  #17](https://github.com/sensu/uchiwa-web/pull/17)

#### Other
- Hide 'Issue Request' button for standalone checks -
[Issue #754](https://github.com/sensu/uchiwa/issues/754)
- Replace 'never' expiration with 'on resolve', if applicable, on silenced view -
[Issue #764](https://github.com/sensu/uchiwa/issues/764)

### 1.1.3-2 (2018-02-15)
#### Bug Fixes
- Fix syntax error in RPM post script - 
[Issue #760](https://github.com/sensu/uchiwa/issues/760)

### 1.1.3 (2018-01-30)
#### Bug Fixes
- Fix typo in TLS cipher list - 
[PR #759](https://github.com/sensu/uchiwa/pull/759)

### 1.1.2 (2018-01-26)
#### Bug Fixes
- Validate both check & client subscriptions when determing if a check is
silenced - [Issue #755](https://github.com/sensu/uchiwa/issues/755)
- Fix logrotate configuration file permission - 
[Issue #757](https://github.com/sensu/uchiwa/issues/757)
- Fix Travis CI builds for both uchiwa & uchiwa-web projects
- Upgrade `istanbul` & `grunt-contrib-watch` dependencies

### 1.1.1 (2018-01-10)
#### Bug Fixes
- Fix the display of Uchiwa configuration in info view  -
[PR #753](https://github.com/sensu/uchiwa/pull/753)
- Fix the Uchiwa logo in Firefox -
[Issue #749](https://github.com/sensu/uchiwa/issues/749)

#### Other
- Allow custom build of packages with Docker -
[PR #718](https://github.com/sensu/uchiwa/pull/718)

### 1.1.0 (2017-12-15)
#### Features
- Removed 3DES cipher algorithms from default suite -
[Issue #688](https://github.com/sensu/uchiwa/issues/688)
- The TLS config is now configurable for HTTPS server -
[Issue #688](https://github.com/sensu/uchiwa/issues/688)

### Bug Fixes
- Fixed the automated Docker images build process -
[Issue #742](https://github.com/sensu/uchiwa/issues/742)

### 1.0.0 (2017-11-30)
#### Bug Fixes
- Docker images now use Go 1.9.2 -
[Issue #731](https://github.com/sensu/uchiwa/issues/731)
- Respect line breaks of results output -
[Issue #708](https://github.com/sensu/uchiwa/issues/708)
- Prevent a panic when the response content length is unknown -
[Issue #725](https://github.com/sensu/uchiwa/issues/725)

#### Other
- Packages now provide a logrotate configuration -
[Issue #703](https://github.com/sensu/uchiwa/issues/703)

### 0.26.3 (2017-10-05)
#### Other
- Use pagination on Sensu API events endpoint -
[Commit 2ba0508b](https://github.com/sensu/uchiwa/commit/2ba0508b3bb95426274d982d03949c5c551d2fcc)
- Add tracing attribute to enable HTTP client tracing -
[PR #720](https://github.com/sensu/uchiwa/pull/720)

### 0.26.2 (2017-10-02)
#### Bug Fixes
- Properly initialize Sensu APIs HTTP clients

### 0.26.1 (2017-09-29)
#### Other
- Optimize the HTTP client used with Sensu APIs -
[PR #719](https://github.com/sensu/uchiwa/pull/719)

### 0.26.0 (2017-09-22)
#### Features
- Additional Sensu server metadata in datacenters view -
[Issue #706](https://github.com/sensu/uchiwa/issues/706)
- Add a type filter to the checks view -
[Issue #310](https://github.com/sensu/uchiwa/issues/310)
- Support issuing check request for a particular client -
[PR #2](https://github.com/sensu/uchiwa-web/pull/2)

#### Bug Fixes
- Better handling of Sensu objects starting with a `.` -
[Issue #710](https://github.com/sensu/uchiwa/issues/710)

#### Other
- Improve event status filtering by handling unknown statuses greater than 3
[PR #1](https://github.com/sensu/uchiwa-web/pull/1)
- New logo!

### 0.25.3 (2017-07-25)
#### Bug Fixes
- Consider clients subscriptions when determining if a check is silenced -
[Issue #689](https://github.com/sensu/uchiwa/issues/689)
- Fix URL shortening with special characters -
[Issue #690](https://github.com/sensu/uchiwa/issues/690)
- Fix the client status filter -
[Commit 0a66d72](https://github.com/sensu/uchiwa-web/commit/0a66d72930560b2050f573dab783100f4661cfb4)

#### Other
- Display the silenced creation timestamp -
[Commit 47db92f](https://github.com/sensu/uchiwa-web/commit/47db92f00102a500c0e26bb52d21460b94ff4414)

### 0.25.2 (2017-05-16)
#### Bug Fixes
- Prevent empty events view upon login
- Events output must be truncated with ellipsis

### 0.25.1 (2017-05-10)
#### Bug Fixes
- Fix checkboxes in the checks view -
[Issue #681](https://github.com/sensu/uchiwa/issues/681)

### 0.25.0 (2017-05-05)
#### Features
- Client creation is now possible from the clients view -
[PR #678](https://github.com/sensu/uchiwa/pull/678)
- Client edition is now possible for all clients -
[PR #678](https://github.com/sensu/uchiwa/pull/678)
- Identify flapping checks in the events view -
[Issue #673](https://github.com/sensu/uchiwa/issues/673)

#### Bug Fixes
- Properly wrap text with URLs in the client view -
[Issue #679](https://github.com/sensu/uchiwa/issues/679)
- Fix bulk silencing for clients with similar names -
[Issue #677](https://github.com/sensu/uchiwa/issues/677)

#### Other
- Add protection against Clickjack/UI Redress attacks -
[PR #680](https://github.com/sensu/uchiwa/pull/680)

### 0.24.0 (2017-04-13)
#### Features
- Add a check view - [Issue #672](https://github.com/sensu/uchiwa/issues/672)
- Add a 'not healthy' option to status filter in clients view -
[Issue #445](https://github.com/sensu/uchiwa/issues/445)

#### Bug Fixes
- Fix regex queries on fields within the check or client object -
[Issue #674](https://github.com/sensu/uchiwa/issues/674)

#### Other
- Shorten long URLs in the client view -
[Issue #629](https://github.com/sensu/uchiwa/issues/629)
- Prevent a user from accessing the login page if authentication is disabled
- Force the browser location to the login page upon logout
- Replace Godeps with govendor for Go vendoring -
[Commit 91fd100](https://github.com/sensu/uchiwa/commit/91fd1009bf207ee8ffbf674f88cc6491a2b186d6)

### 0.23.1 (2017-03-29)
#### Bug Fixes
- Prevent caching of deprecated assets when upgrading -
[Issue #666](https://github.com/sensu/uchiwa/issues/666)
- Fix bulk silencing when choosing a duration and/or a reason -
[Issue #667](https://github.com/sensu/uchiwa/issues/667)
- Support negative lookahead for search queries -
[Issue #664](https://github.com/sensu/uchiwa/issues/664)

#### Other
- Increase verbosity for datacenters connection errors -
[Issue #663](https://github.com/sensu/uchiwa/issues/663)
- Upgrade to Angular 1.6 -
[Issue #653](https://github.com/sensu/uchiwa/issues/653)
- Improve responsive design with long text -
[Issue #641](https://github.com/sensu/uchiwa/issues/641)

### 0.23.0 (2017-03-24)
#### Bug Fixes
- Fix a visual glitch when refreshing a page on Firefox -
[Issue #654](https://github.com/sensu/uchiwa/issues/654)

#### Features
- Refactoring of the silencing modal window -
[Issue #644](https://github.com/sensu/uchiwa/issues/644)
- Support edition of proxy clients in the Sensu client registry -
[Issue #651](https://github.com/sensu/uchiwa/issues/651)
- Display new information about Sensu servers, available since Sensu 0.28, in a new
datacenter view - [Issue #652](https://github.com/sensu/uchiwa/issues/652)
- Support custom favicons - [Issue #606](https://github.com/sensu/uchiwa/issues/606)

#### Other
- Set the Secure flag on authentication cookies when serving over HTTPS -
[Issue #656](https://github.com/sensu/uchiwa/issues/656)

### 0.22.2 (2017-03-10)
#### Bug Fixes
- Add Cross Site Request Forgery (XSRF) protection - [Issue #638](https://github.com/sensu/uchiwa/issues/638)

#### Other
- Use Travis CI for the build pipeline - [Issue #633](https://github.com/sensu/uchiwa/issues/633)

### 0.22.1 (2017-02-06)
#### Bug Fixes
- The bulk actions menu should be shown when at least one element is selected

#### Other
- Default to the new option _On resolve_ when silencing an element - [Issue #631](https://github.com/sensu/uchiwa/issues/631)

### 0.22.0 (2017-01-25)
#### Features
- Responsive web design - [Issue #616](https://github.com/sensu/uchiwa/issues/616)

#### Bug Fixes
- Prevent reflected cross-site scripting (XSS) attacks - [Issue #618](https://github.com/sensu/uchiwa/issues/618)
- Prevent race conditions in the logger package - [Issue #621](https://github.com/sensu/uchiwa/issues/621)
- Fix recursive searches while using key-value queries - [Issue #624](https://github.com/sensu/uchiwa/issues/624)


### 0.21.0 (2016-12-08)
#### Features
- Increase accuracy of the silenced attribute for Sensu events - [Issue #602](https://github.com/sensu/uchiwa/issues/602)
- Add the TTL attribute to the checks view - [Issue #496](https://github.com/sensu/uchiwa/issues/496)

#### Bug Fixes
- Use relative URLs for Uchiwa API resources - [Issue #600](https://github.com/sensu/uchiwa/issues/600)
- Fix silencing from the client view - [Issue #601](https://github.com/sensu/uchiwa/issues/601)
- Prevent an empty popover when all clients are healthy - [Issue #607](https://github.com/sensu/uchiwa/issues/607)
- Display the output of all checks on the client view - [Issue #612](https://github.com/sensu/uchiwa/issues/612)

### 0.20.2 (2016-11-21)
#### Bug Fixes
- Use the ngSrc directive to prevent failed requests when loading the logo URL -
[Issue #599](https://github.com/sensu/uchiwa/issues/599)
- Also display the custom logo on the login page
- Fix the user dropdown when authenticated
- Fix the theme selection on the settings view
- Adjust the error description if the Uchiwa backend config is unavailable -
[Commit 5de458a](https://github.com/sensu/uchiwa/commit/5de458a)

### 0.20.1 (2016-11-17)
#### Bug Fixes
- Prevent any caching of the index.html file in order to facilitate the upgrade
process - [PR #597](https://github.com/sensu/uchiwa/pull/597)
- Do not apply a filter on the status attribute in the client view
- Fix a typographical error in the clients popover of the sidebar

### 0.20.0 (2016-11-14)
#### Features
- Added various users-level configuration attributes to customize Uchiwa
- Support regex and key:value search queries
- Sidebar popovers redesign
- Aggregates can now be deleted
- Show the reason in the silenced view

#### Other
- Refactoring of the Angular services
- The checks and subscriptions filters now only show values associated with a datacenter if one is selected in the datacenters filter
- Improve notification messages - [Issue #467](https://github.com/sensu/uchiwa/issues/467)
- Move iframes to their own panels in the client view - [Issue #360](https://github.com/sensu/uchiwa/issues/360)
- Pre-0.26 Sensu APIs are now marked as deprecated in the logs - [Issue #590](https://github.com/sensu/uchiwa/issues/590)
- Refactoring of the Angular bootstrapping
- Remove underscore.js dependency

### 0.19.0 (2016-10-16)
#### Features
- Allow silencing of checks and subscriptions across datacenters
- Display rich aggregates data

#### Bug Fixes
- The pagination counter should consider the filters applied - [Issue #431](https://github.com/sensu/uchiwa/issues/431)
- Do not panic when an encrypted password is invalid - [Issue #549](https://github.com/sensu/uchiwa/issues/549)

### 0.18.2 (2016-09-19)
#### Bug Fixes
- Fix the silenced filters - [Issue #565](https://github.com/sensu/uchiwa/issues/565)  
- Do not move an image from the command attribute to its own box - [Issue  #558](https://github.com/sensu/uchiwa/issues/558)

#### Other
- Allow choice of silencing entries when un-silencing an item
- Allow choice of custom expiration when creating a silence entry - [Issue #570](https://github.com/sensu/uchiwa/issues/570)  
- The /health endpoint can return a 503 HTTP response code on error - [Issue #557](https://github.com/sensu/uchiwa/issues/557)

### 0.18.1 (2016-09-09)
#### Bug Fixes
- Fix silencing with no expiration

### 0.18.0 (2016-09-07)
**Requires Sensu >= 0.26**

#### Bug Fixes
- Prevent race condition when using the Uchiwa filters - [Issue #543](https://github.com/sensu/uchiwa/issues/543)
- Fix styling of the uchiwa-dark theme

#### Features
- Support built-in silencing in Sensu 0.26 - [Issue #539](https://github.com/sensu/uchiwa/issues/539)
- Filter per-client subscriptions - [Issue #534](https://github.com/sensu/uchiwa/issues/534)
- Add support for encrypted passwords - [PR #448](https://github.com/sensu/uchiwa/issues/448)
- Display last_ok attribute in events view

### 0.17.1 (2016-08-02)
#### Bug Fixes
- Remove various debugging traces - [Commit #d249aa4](https://github.com/sensu/uchiwa/commit/d249aa4)

#### Other
- Upgrade to Go 1.6.3 with vendoring support - [PR #528](https://github.com/sensu/uchiwa/pull/528)
- The filters package now implements an interface for easier use - [PR #528](https://github.com/sensu/uchiwa/pull/528)
- Refactoring of the authorization logic - [PR #528](https://github.com/sensu/uchiwa/pull/528)

### 0.17.0 (2016-07-20)
#### Features
- Add API token-based authentication - [PR #525](https://github.com/sensu/uchiwa/pull/525)

#### Bug Fixes
- Prevent old alerts to be displayed inadvertently - [Issue #512](https://github.com/sensu/uchiwa/issues/512)
- Fix iframes support - [Issue #508](https://github.com/sensu/uchiwa/issues/508)

#### Other
Use Alpine Linux as the base image for Docker images - [PR #498](https://github.com/sensu/uchiwa/pull/498)

### 0.16.0 (2016-06-23)
#### Bug Fixes
- The stashes could not be deleted from the stashes view - [Issue #503](https://github.com/sensu/uchiwa/issues/507)
- Incorrect client attributes could be displayed on a client view - [Issue #471](https://github.com/sensu/uchiwa/issues/471)
- The 'Show All' option should not use the current elements number - [Issue #466](https://github.com/sensu/uchiwa/issues/466)
- The relative timestamp was not properly calculated on a stash view - [Issue #456](https://github.com/sensu/uchiwa/issues/456)

#### Other
- Refactoring of the Uchiwa API endpoints - [PR #513](https://github.com/sensu/uchiwa/pull/513)

### 0.15.0 (2016-06-02)
#### Features
- Add support for upcoming Sensu 0.24.0 - [PR #500](https://github.com/sensu/uchiwa/pull/500)

### 0.14.5 (2016-05-10)
#### Bug Fixes
- Fix i386 packages - [PR #31](https://github.com/sensu/uchiwa-build/pull/31)

### 0.14.4 (2016-05-10)
#### Bug Fixes
- Add missing syntax highlighting on the client view - [1089a7d](https://github.com/sensu/uchiwa-web/commit/1089a7d75ca8810c31fc44492acb07dc402faa5a)
- Fix the checks filter on the events view
- Prevent infinite pagination loop on the Sensu API - [PR #478](https://github.com/sensu/uchiwa/pull/478)

### 0.14.3 (2016-03-03)
#### Features
- Add a detailed stash view - [PR #90](https://github.com/sensu/uchiwa-web/pull/90)

#### Bug Fixes
- The HTML code for syntax highlighting was not always properly processed

### 0.14.2 (2016-02-02)
#### Features
- Add support for serving content over HTTPS - [PR #441](https://github.com/sensu/uchiwa/pull/441)

### 0.14.1 (2016-01-15)
#### Bug Fixes
- Fix version number for Debian packages - [Issue #26](https://github.com/sensu/uchiwa-build/issues/26)

#### Other
- Improved logging with log levels - [PR #425](https://github.com/sensu/uchiwa/pull/425)
- Upgrade to Go 1.5.3

### 0.14.0 (2016-01-05)
#### Features
- Load Uchiwa configuration from directories - [PR #416](https://github.com/sensu/uchiwa/pull/416)
- Issue check requests from the checks view - [Issue #141](https://github.com/sensu/uchiwa/issues/141)
- Delete a client check result, **requires Sensu 0.21.0 or later** - [PR #419](https://github.com/sensu/uchiwa/pull/419)

### 0.13.0 (2015-11-22)
#### Features
- Datacenters high availability (support multiple APIs for the same datacenter) - [Issue #173](https://github.com/sensu/uchiwa/issues/173) - [Docs](http://docs.uchiwa.io/en/latest/configuration/sensu/#datacenters-high-availability)
- Static RSA keys for the JSON Web Tokens signature - [Issue #394](https://github.com/sensu/uchiwa/issues/394) - [Docs](http://docs.uchiwa.io/en/latest/configuration/uchiwa/#static-rsa-keys)
- Upgrade angular-toastr version to 1.6.0

#### Bug Fixes
- Fix the _Hide Silenced Clients_ filter on the events view - [Issue #412](https://github.com/sensu/uchiwa/issues/412)
- Display an error message when Uchiwa fails to contact its backend API
- Make sure to update the health and metrics data on all views
- Tweak the badges position on the sidebar

### 0.12.1 (2015-11-05)
#### Bug Fixes
- History for all checks was not properly displayed on the client view - [Issue #404](https://github.com/sensu/uchiwa/issues/404)
- Better handling of JIT clients with no timestamp

### 0.12.0 (2015-10-20)
#### Features
- Major performance improvements on the frontend, especially when manipulating ten of thousands of elements - [Issue #399](https://github.com/sensu/uchiwa/issues/399)
- Use pagination when querying the Sensu API - [Issue #397](https://github.com/sensu/uchiwa/issues/397)
- Refactoring of the client history in order to display rich information for all checks, including standalones, on the client view. Deprecating support for Sensu 0.12. - [Issue #395](https://github.com/sensu/uchiwa/issues/395)

### 0.11.2 (2015-10-04)
#### Bug Fixes
- Prevent undefined object when evaluating scope.metrics object in SidebarController - [Issue #387](https://github.com/sensu/uchiwa/issues/387)
- Simplify Dockerfile and upgrade golang docker image to 1.5.1 - [Issue #391](https://github.com/sensu/uchiwa/issues/391)

### 0.11.1 (2015-09-23)
#### Bug Fixes
- Redirection to the login page should remove all query strings - [Issue #385](https://github.com/sensu/uchiwa/issues/385)
- Add versioning to JS & CSS files to avoid caching with upgrades - [Issue #386](https://github.com/sensu/uchiwa/issues/386)
- Set *success* style to the events sidebar icon when we have no events
- Prevent errors when an API endpoint returns null
- Dependency cleanup

### 0.11.0 (2015-09-22)
#### Features
- Implement a RESTful API and remove the *get_sensu* endpoint for Uchiwa backend - [Issue #378](https://github.com/sensu/uchiwa/pull/378)
- Major frontend performance and stability improvement: use the newer Uchiwa RESTful API and store data into $scope instead of $rootScope
- Allow bulk removal of stashes
- Also display client's images in a dedicated panels on the client view - [Issue #361](https://github.com/sensu/uchiwa/issues/361)
- Add progress bar into aggregate view
- Upgrade to Go 1.5.1

#### Bug Fixes
- Allow text selection without immediately firing ng-click - [Issue #262](https://github.com/sensu/uchiwa/issues/262)
- Break long datacenter name into multiple lines - [Issue #368](https://github.com/sensu/uchiwa/issues/368)
- Add username to stash content - [Issue #356](https://github.com/sensu/uchiwa/issues/356)
- The /results Sensu API endpoint is not required yet - [Issue #379](https://github.com/sensu/uchiwa/pull/379)
- Two events with the same client and check names could be mixed - [Issue #375](https://github.com/sensu/uchiwa/issues/375)
- Recover from an unexpected type assertion when processing a client
- Prevent multiple status code within a single HTTP response on the /health endpoint
- Properly display username in navbar if authentication is enabled

### 0.10.4 (2015-09-01)
#### Bug Fixes
- Order alphabetically the items in the checks filter
- Fix client view table for RO users
- Visual improvements to client view when resizing to a narrow view
- Fix favicon for Firefox - [Issue #376](https://github.com/sensu/uchiwa/pull/367)
- Add support for Sensu Enteprise OpenLDAP driver - [Issue #369](https://github.com/sensu/uchiwa/pull/369)
- Add support for Sensu Enteprise audit logging - [Issue #370](https://github.com/sensu/uchiwa/pull/370)

### 0.10.3 (2015-08-03)
#### Features
- Add status filter on the clients and events views
- Add support for Sensu Enterprise dashboard

### 0.10.2 (2015-07-23)
#### Features
- Add username to stash content - [Issue #356](https://github.com/sensu/uchiwa/issues/356)
- Replace silenced and critical icons
- Add support for audit logging - Sensu Enterprise Dashboard

#### Bug Fixes
- Reimplement iframe support - [Issue #354](https://github.com/sensu/uchiwa/issues/354)
- Display any error with the http.ListenAndServe method - [Issue #352](https://github.com/sensu/uchiwa/issues/352)

### 0.10.1 (2015-06-30)
#### Bug Fixes
- Fix check result view for standalone checks - [Issue #350](https://github.com/sensu/uchiwa/issues/350)

### 0.10.0 (2015-06-29)

#### Features
- Multiple users (RO & RW) can be defined in the configuration - [Issue #343](https://github.com/sensu/uchiwa/pull/343)
- The theme setting is now saved in a cookie - [Issue #331](https://github.com/sensu/uchiwa/issues/331)
- Display the output of all checks in the check result view - [Issue #346](https://github.com/sensu/uchiwa/issues/346)
- Enhancements to the check result images

#### Bug Fixes
- The info view might have been incomplete
- Disable autocapitalization and autocorrection on the login view - [Issue #296](https://github.com/sensu/uchiwa/issues/296)
- Remove unsupported characters in datacenter name - [Issue #279](https://github.com/sensu/uchiwa/issues/279)
- Continue to pull the client details even when an error is returned - [Issue #265](https://github.com/sensu/uchiwa/issues/265)

### 0.9.1 (2015-06-10)
#### Bug Fixes
- Performance issues - [Issue #337](https://github.com/sensu/uchiwa/issues/337)

### 0.9.0 (2015-06-09)
#### Features
- Display the output for all checks - [Issue #322](https://github.com/sensu/uchiwa/issues/322)
- Various fixes and improvements to the backend - [Issue #330](https://github.com/sensu/uchiwa/pull/330)
  - Godep is now used to manage vendored dependencies
  - Support for Sensu Enterprise dashboard features
  - Refactoring of the Go packages
- Allow filtering by check on the checks view
- Include a result count when searching

#### Bug Fixes
- Better handling of invalid events - [Issue #332](https://github.com/sensu/uchiwa/issues/332)
- A stash can only start now and not in the future

### 0.8.1 (2015-05-05)
#### Features
- Add profile picture to the navbar when logged

#### Bug Fixes
- Allow stash creation with no expiration - [Issue #319](https://github.com/sensu/uchiwa/issues/319)
- Upgrade to angular-bootstrap 0.13.0 - [Issue #319](https://github.com/sensu/uchiwa/issues/319)

### 0.8.0 (2015-04-29)
#### Features
- Import the palourde/auth library within the Uchiwa repository - [Issue #314](https://github.com/sensu/uchiwa/pull/314)
- Refactoring of the stashes API endpoints - [Issue #317](https://github.com/sensu/uchiwa/pull/317)
- Add relative times to stashes and clients views
- Add support for Github authentication driver (Sensu Enterprise)
- Add support for LDAP authentication driver (Sensu Enterprise)

#### Bug Fixes
- Allow stash creation with expiration longer than a few days - [Issue #301](https://github.com/sensu/uchiwa/issues/301)
- Datacenter filter now performs a strict comparison - [Issue #307](https://github.com/sensu/uchiwa/issues/307)
- Resolved events are now cleared from the clients view - [Issue #309](https://github.com/sensu/uchiwa/issues/309)

### 0.7.1 (2015-04-01)
#### Features
- Show relative times for events
- Add datacenters view
- Order events by status, then by most recent

#### Bug Fixes
- Remove *standalone* property from checks view
- Show logout button when authentication is enabled
- Ship Google fonts with Uchiwa

### 0.7.0 (2015-03-13)
#### Features
- Reorganize the navbar and sidebar
- Panels styling
- Improvements to the aggregates view

#### Bug Fixes
- Unselect events after action in events view
- Fix stash expiration timestamp in stashes view

### 0.6.0 (2015-02-26)
#### Features
- Upgrade to AngularJS 1.3 - [Issue #160](https://github.com/sensu/uchiwa/issues/160)
- Redesign of the panels header
- Add tooltip of items subscriptions on clients and events views

### 0.5.1 (2015-02-21)

#### Features
- Add *Silenced Clients* option to the Hide menu and show the acknowledgment status of the client in the events view

#### Bug Fixes
- Register *uchiwa-web* as a bower package to prevent dependencies issues - [Issue #272](https://github.com/sensu/uchiwa/issues/272)
- Fix build issue with 0.5.0 release - [Issue #273](https://github.com/sensu/uchiwa/issues/273)

### 0.5.0 (2015-02-14)

#### Features
- Add custom date functionality for stashes - [Issue #251](https://github.com/sensu/uchiwa/pull/251)
- Aggregates support - [Issue #255](https://github.com/sensu/uchiwa/pull/255)
- Use JSON Web Tokens (JWT) instead of HTTP basic authentication and implement a login page

#### Bug Fixes
- Catch possible exception while asserting the check name - [Issue #259](https://github.com/sensu/uchiwa/pull/259)
- Display images when an URL has a port number
- Hide events with no client and no check
- Client names may not be unique

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
