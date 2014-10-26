'use strict';

module.exports = function (grunt) {
  require('load-grunt-tasks')(grunt);
  require('time-grunt')(grunt);

  grunt.initConfig({
    settings: {
      public: 'public'
    },
    jshint: {
      options: {
        jshintrc: '.jshintrc',
        reporter: require('jshint-stylish')
      },
      all: [
        'Gruntfile.js',
        '<%= settings.public %>/js/{,*/}*.js',
      ]
    },
    karma: {
      unit: {
        configFile: 'test/karma/conf.js'
      }
    },
    lintspaces: {
      all: {
        src: [
          'public/js/**/*.js'
        ],
        options: {
          newline: true,
          newlineMaximum: 2,
          indentation: 'spaces',
          spaces: 2,
          trailingspaces: true
        }
      }
    },
    sass: {
      dist: {
        files: [{
          expand: true,
          cwd: '<%= settings.public %>/css',
          src: ['**/*.scss'],
          dest: '<%= settings.public %>/css',
          ext: '.css'
        }]
      }
    },
    watch: {
      scss: {
        files: ['<%= settings.public %>/css/**/*.scss'],
        tasks: ['sass:dist'],
        options: {
          spawn: false
        }
      }
    }
  });

  grunt.registerTask('dev', [
    'watch:scss'
  ]);

  grunt.registerTask('lint', [
    'newer:jshint'
  ]);

  grunt.registerTask('default', [
    'sass',
    'jshint',
    'lintspaces',
    'karma:unit'
  ]);
};
