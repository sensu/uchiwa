'use strict';

describe('directives', function() {
  var scope;
  var element;

  beforeEach(module('uchiwa'));
  beforeEach(inject(function() {
    scope = jasmine.createSpyObj('scope', ['$on']);
    element = jasmine.createSpyObj('element', ['tooltip', 'attr']);
  }));

  describe('bootstrapTooltip', function() {

    it('should be restricted to elements and attributes', inject(function(bootstrapTooltipDirective) {
      expect(bootstrapTooltipDirective[0].restrict).toBe('EA');
    }));

    it('should have a link method', inject(function(bootstrapTooltipDirective) {
      expect(bootstrapTooltipDirective[0].link).toBeDefined();
    }));

    it('should call element.tooltip() when calling link', inject(function(bootstrapTooltipDirective) {
      bootstrapTooltipDirective[0].link({}, element);
      expect(element.tooltip).toHaveBeenCalled();
    }));

  });

  describe('siteTheme', function() {

    it('should be restricted to elements and attributes', inject(function(siteThemeDirective) {
      expect(siteThemeDirective[0].restrict).toBe('EA');
    }));

    it('should have a link method', inject(function(siteThemeDirective) {
      expect(siteThemeDirective[0].link).toBeDefined();
    }));

    it('should define themes', inject(function(siteThemeDirective) {
      siteThemeDirective[0].link(scope, element);
      expect(scope.themes).toBeDefined();
      expect(scope.themes.length).toBeGreaterThan(0);
    }));

    it('should listen for theme:changed event', inject(function(siteThemeDirective) {
      siteThemeDirective[0].link(scope, element);
      expect(scope.$on).toHaveBeenCalledWith('theme:changed', jasmine.any(Function));
    }));
  });
});