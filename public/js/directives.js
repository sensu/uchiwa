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