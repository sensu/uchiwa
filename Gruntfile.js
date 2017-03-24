'use strict';

module.exports = function (grunt) {
  grunt.loadNpmTasks("grunt-cache-breaker");

  grunt.initConfig({
    cachebreaker: {
      dev: {
        options: {
          match: [
            'app.js',
            'bootstrap.js',
            'common.js',
            'constants.js',
            'controllers.js',
            'directives.js',
            'factories.js',
            'filters.js',
            'services.js',
            'angular.min.js',
            'angular-cookies.min.js',
            'angular-resource.min.js',
            'angular-route.min.js',
            'angular-sanitize.min.js',
            'async.js',
            'angular-toastr.tpls.min.js',
            'ui-bootstrap-tpls.min.js',
            'highlight.pack.js',
            'moment.min.js',
            'angular-moment.min.js',
            'angular-moment-picker.min.js',
            'angular-gravatar.min.js',
            'jsoneditor.min.js',
            'ng-jsoneditor.js',
            'bootstrap-uchiwa.css',
            'font-awesome.min.css',
            'angular-toastr.min.css',
            'tomorrow.css',
            'angular-moment-picker.min.css',
            'jsoneditor.min.css',
            'uchiwa-default.css',
            'ua-parser.min.js'
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
