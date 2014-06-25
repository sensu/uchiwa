'use strict';

module.exports = function (grunt) {
  require('load-grunt-tasks')(grunt);
  require('time-grunt')(grunt);

  grunt.initConfig({
    settings: {
      entryPoint: 'app.js',
      lib: 'lib',
      public: 'public'
    },
    jshint: {
      options: {
        jshintrc: '.jshintrc',
        reporter: require('jshint-stylish')
      },
      all: [
        'Gruntfile.js',
        '<%= settings.entryPoint %>',
        '<%= settings.lib %>/{,*/}*.js',
        '<%= settings.public %>/js/{,*/}*.js',
      ]
    }
  });

  grunt.registerTask('default', [
    'newer:jshint'
  ]);
};
