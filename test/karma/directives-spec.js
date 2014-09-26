'use strict';

describe('directives', function () {
  var scope;
  var element;

  beforeEach(module('uchiwa'));
  beforeEach(inject(function () {
    scope = jasmine.createSpyObj('scope', ['$on', '$watch']);
    element = jasmine.createSpyObj('element', ['tooltip', 'attr', 'removeAttr', 'addClass']);
  }));

  describe('siteTheme', function () {

    it('should be restricted to elements and attributes', inject(function (siteThemeDirective) {
      expect(siteThemeDirective[0].restrict).toBe('EA');
    }));

    it('should have a link method', inject(function (siteThemeDirective) {
      expect(siteThemeDirective[0].link).toBeDefined();
    }));

    it('should define themes', inject(function (siteThemeDirective) {
      siteThemeDirective[0].link(scope, element);
      expect(scope.themes).toBeDefined();
      expect(scope.themes.length).toBeGreaterThan(0);
    }));

    it('should listen for theme:changed event', inject(function (siteThemeDirective) {
      siteThemeDirective[0].link(scope, element);
      expect(scope.$on).toHaveBeenCalledWith('theme:changed', jasmine.any(Function));
    }));
  });

  describe('statusGlyph', function () {

    it('should be restricted to elements and attributes', inject(function (statusGlyphDirective) {
      expect(statusGlyphDirective[0].restrict).toBe('EA');
    }));

    it('should have a link method', inject(function (statusGlyphDirective) {
      expect(statusGlyphDirective[0].link).toBeDefined();
    }));

    it('should add classes when calling link', inject(function (statusGlyphDirective) {
      var attrs = {
        statusGlyph: 'client.style'
      };
      statusGlyphDirective[0].link(scope, element, attrs);

      expect(scope.$watch).toHaveBeenCalledWith(attrs.statusGlyph, jasmine.any(Function));
    }));

  });

});
