# gosensu

Gosensu is an golang wrapper around the Sensu API. It's a convenient wrapper to make REST calls
  to the Sensu API. It is used by the Golang version of the backend for [Uchiwa](http://github.com/sensu/uchiwa).

## Usage

A quick example:

```golang
sensu := New("My sensu API", "", "http://your.SENSU_SERVER_URL.tld"), 15, "username", "secret")
events = sensu.GetEvents()
```

See the included godoc documentation for more information about usage.

## Caveats

This library is the first aim of the authors at Golang. Thank you for any feedback for improvement
on styles or golang idioms!

## Contributing

Everyone is welcome to submit patches. Whether your pull request is a bug fix or introduces new classes or functions to the project, we kindly ask that you include tests for your changes. Even if it's just a small improvement, a test is necessary to ensure the bug is never re-introduced.

### Testing

Some tests are included with this project (remember, the author is still a noob ;)). Testing
an API wrapper is always somewhat a challenge since mocking everything can be challenging and having
a full blown client cumbersome.

A quick testing can be done by running 'make canned' which will test against a [Canned](https://github.com/palourde/sensu-canned)
powered mockup of the API (driven by static JSON files, via Node's canned), running in an Heroku setup.

More "real" testing can be done by running "make test" that will be runned against a local installation
of sensu (currenctly hardcoded to localhost port 8889). This code have been developped against
a VM built so support Uchiwa's developpement that can be found at https://github.com/palourde/uchiwa-sensu

We recommend to run these tests (or even better, improve them!) before submitting a Pull Request.

## Authors
* Author: [Benoit Caron][author] (<benoit@patentemoi.ca>)
* Contributor: [Simon Plourde] (<simon.plourde@gmail.com>)

## License
MIT (see [LICENSE][license])

[author]:                 https://github.com/bencaron
[license]:                https://github.com/bencaron/gosensu/blob/master/LICENSE
