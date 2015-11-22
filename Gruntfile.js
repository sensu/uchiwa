'use strict';

module.exports = function (grunt) {
  grunt.loadNpmTasks("grunt-cache-breaker");

  grunt.initConfig({
    cachebreaker: {
      dev: {
        options: {
          match: [
            'app.js',
            'constants.js',
            'controllers.js',
            'directives.js',
            'factories.js',
            'filters.js',
            'providers.js',
            'services.js',
            'angular.min.js',
            'angular-cookies.min.js',
            'angular-route.min.js',
            'angular-sanitize.min.js',
            'async.js',
            'underscore.js',
            'angular-toastr.tpls.min.js',
            'ui-bootstrap-tpls.min.js',
            'highlight.pack.js',
            'moment.min.js',
            'angular-moment.min.js',
            'angular-gravatar.min.js',
            'bootstrap-uchiwa.css',
            'font-awesome.min.css',
            'angular-toastr.min.css',
            'tomorrow.css',
            'uchiwa-default.css'
          ],
          replacement: 'time'
        },
        files: {
          src: ['public/index.html']
        }
      }
    }
  });

  grunt.registerTask('default', [
    'cachebreaker'
  ]);
};
