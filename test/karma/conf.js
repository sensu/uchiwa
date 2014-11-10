// karma.conf.js
module.exports = function(config) {
  config.set({
    basePath : '../../',
    frameworks: ['jasmine'],
    files : [
      'public/bower_components/underscore/underscore.js',
      'public/bower_components/angular/angular.js',
      'public/bower_components/angular-cookies/angular-cookies.js',
      'public/bower_components/angular-route/angular-route.js',
      'public/bower_components/angular-sanitize/angular-sanitize.js',
      'public/bower_components/angular-mocks/angular-mocks.js',
      'public/bower_components/angular-toastr/dist/angular-toastr.min.js',
      'public/bower_components/angular-bootstrap/ui-bootstrap-tpls.min.js',
      'public/bower_components/moment/min/moment.min.js',
      'public/bower_components/highlightjs/highlight.pack.js',
      'public/js/**/*.js',
      'test/karma/**/*.js'
    ],
    reporters: ['junit', 'coverage', 'dots'],
    coverageReporter: {
      type: 'html',
      dir: 'build/coverage/'
    },
    preprocessors: {
      'public/js/**/*.js': ['coverage']
    },
    junitReporter: {
      outputFile: 'build/karma/test-results.xml'
    },
    port: 8876,
    colors: true,
    logLevel: config.LOG_INFO,
    autoWatch: true,
    browsers: ['PhantomJS'],
    captureTimeout: 60000,
    singleRun: true
  });
};
