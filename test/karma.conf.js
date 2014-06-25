// karma.conf.js
module.exports = function(config) {
  config.set({
    basePath : '../',
    frameworks: ['jasmine'],
    files : [
      'public/bower_components/jquery/dist/jquery.js',
      'public/bower_components/angular/angular.js',
      'public/bower_components/angular-cookies/angular-cookies.js',
      'public/bower_components/angular-route/angular-route.js',
      'public/bower_components/angular-mocks/angular-mocks.js',
      'public/bower_components/angular-socket.io-mock/angular-socket.io-mock.js',
      'public/bower_components/toastr/toastr.min.js',
      'public/js/**/*.js',
      'test/unit/**/*.js'
    ],
    reporters: ['junit', 'coverage'],
    coverageReporter: {
      type: 'html',
      dir: 'build/coverage/'
    },
    preprocessors: {
      'public/js/**/*.js': ['coverage']
    },
    junitReporter: {
      outputFile: 'build/unit/test-results.xml'
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