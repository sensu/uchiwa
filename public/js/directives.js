'use strict';

var directiveModule = angular.module('uchiwa.directives', []);

directiveModule.directive('siteTheme', ['$cookieStore', 'settings', function ($cookieStore, settings) {
  return {
    restrict: 'EA',
    link: function (scope, element) {
      scope.themes = [
        {
          name: 'default'
        },
        {
          name: 'dark'
        }
      ];
      var lookupTheme = function (themeName) {
        return scope.themes[scope.themes.map(function (t) {
          return t.name;
        }).indexOf(themeName)];
      };
      var setTheme = function (theme) {
        var themeName = angular.isObject(theme) && angular.isDefined(theme.name) ? theme.name : settings.theme;
        scope.currentTheme = lookupTheme(themeName);
        $cookieStore.put('currentTheme', scope.currentTheme);
        var fullThemeName = 'uchiwa-' + scope.currentTheme.name;
        element.attr('href', 'css/' + fullThemeName + '/' + fullThemeName + '.css');
      };
      scope.$on('theme:changed', function (event, theme) {
        setTheme(theme);
      });

      setTheme($cookieStore.get('currentTheme'));
    }
  };
}]);

directiveModule.directive('statusGlyph', ['$filter', function ($filter) {
  return {
    restrict: 'EA',
    link: function(scope, element, attrs) {

      function updateGlyph(style) {
        element.removeAttr('class');
        element.addClass('fa fa-fw');
        switch(style) {
          case 0:
            element.addClass('fa-check-circle');
            break;
          case 1:
            element.addClass('fa-exclamation-circle');
            break;
          case 2:
            element.addClass('fa-bomb');
            break;
          case 3:
            element.addClass('fa-question-circle');
            break;
        }

        var status = $filter('getStatusClass')(style);
        element.addClass('text-' + status);
      }

      scope.$watch(attrs.statusGlyph, function(value) {
        updateGlyph(value);
      });
    }
  };
}]);
