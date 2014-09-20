describe("constants", function () {
  beforeEach(module('uchiwa'));

  var settings;

  describe("settings", function () {
    beforeEach(inject(function ($controller, _$rootScope_, _settings_) {
      settings = _settings_;
    }));

    it("should provide a settings constant", function () {
      expect(angular.isObject(settings)).toBeTruthy();
    });

    it("should provide a theme setting with a value of 'default'", function () {
      var expected = 'default';

      expect(settings.theme).toEqual(expected);
    });
  });
});