'use strict';

describe('directives', function() {
  var scope, element;

  beforeEach(module('ngCookies'));
  beforeEach(module('uchiwa.constants'));
  beforeEach(module('uchiwa.directives'));
  beforeEach(inject(function($rootScope) {
    scope = jasmine.createSpyObj('scope', ['$on']);
    element = jasmine.createSpyObj('element', ['tooltip', 'attr']);
  }));

  describe('bootstrapTooltip', function() {

    it('should be restricted to elements and attributes', inject(function(bootstrapTooltipDirective) {
      expect(bootstrapTooltipDirective[0].restrict).toBe('EA');
    }));

    it('should define link()', inject(function(bootstrapTooltipDirective) {
      expect(bootstrapTooltipDirective[0].link).toBeDefined();
    }));

    it('should call element.tooltip() on link()', inject(function(bootstrapTooltipDirective) {
      bootstrapTooltipDirective[0].link({}, element);
      expect(element.tooltip).toHaveBeenCalled();
    }));

  });

  describe('siteTheme', function() {

    it('should be restricted to elements and attributes', inject(function(siteThemeDirective) {
      expect(siteThemeDirective[0].restrict).toBe('EA');
    }));

    it('should define link()', inject(function(siteThemeDirective) {
      expect(siteThemeDirective[0].link).toBeDefined();
    }));

    it('should create themes property on link()', inject(function(siteThemeDirective) {
      siteThemeDirective[0].link(scope, element);

      expect(scope.themes).toBeDefined();
      expect(scope.themes.length).toBeGreaterThan(0);
      expect(scope.$on).toHaveBeenCalledWith('theme:changed', jasmine.any(Function));
    }));

  });
});