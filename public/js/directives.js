var directiveModule = angular.module('uchiwa.directives', []);

directiveModule.directive('morrisLine', function () {
  return {
    restrict: 'EA',
    replace: true,
    template: '<div class="graph"></div>',
    scope: {
      morrisLine: '='
    },
    link: function (scope, element) {
      scope.morrisLine.element = element;
      var chart = new Morris.Line(scope.morrisLine);
      scope.$watch('morrisLine.data', function (newValue) {
        chart.setData(newValue);
      });
    }
  };
});

directiveModule.directive('bootstrapTooltip', function () {
  return {
    restrict: 'EA',
    link: function (scope, element) {
      element.tooltip();
    }
  };
});

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
        scope.currentTheme = lookupTheme(themeName)
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
