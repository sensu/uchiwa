'use strict';

describe('Uchiwa App', function(){

	describe('Client view', function(){
		beforeEach(function() {
			browser.get('#/');
			element.all(by.css('div.list .client div')).first().click();
		});

		it('should open and close client view modal', function(){
			var modal;

			modal = element(by.id('client-details'));
			expect(modal.isDisplayed()).toBe(true);

			modal.sendKeys(protractor.Key.ESCAPE)
			expect(modal.isDisplayed()).toBe(false);
		});

		it('should display client details', function(){
			expect(element(by.id('client-name')).getText()).toBe('SERVER-0-12-6');
			expect(element(by.id('client-dc')).getText()).toBe('0.12.6');
			expect(element(by.id('client-version')).getText()).toBe('<= 0.12.6-5');
			expect(element(by.id('client-subscriptions')).getText()).toBe('linux');
		});

		it('should list checks', function(){
			var rows = element.all(by.repeater('check in client.history'));
			rows.first().then(function(row){
				var rowElems = row.all(by.tagName('td'));
				rowElems.then(function(cols) {
					expect(cols[0].getText()).toBe('ACTIVE');
					expect(cols[1].getText()).toBe('check_critical');
					expect(cols[2].getText()).toContain('CheckReturn CRITICAL');
					expect(cols[3].getText()).not.toBe('Invalid date');
				});
				expect(row.element(by.id('button-resolve')).getAttribute('class')).not.toContain('btn-disabled');
			});
		});

		it('should display event and check details', function(){
			element.all(by.repeater('check in client.history')).first().click();
			expect(element.all(by.repeater('check in client.history')).get(1).element(by.css('div')).getAttribute('class')).toContain('in');

			element.all(by.repeater('check in client.history')).first().click();
			expect(element.all(by.repeater('check in client.history')).get(1).element(by.css('div')).getAttribute('class')).toContain('collapse');
		});
	});

	describe('Event view', function(){
		beforeEach(function() {
			browser.get('#/');
		});

		it('should display Events as page title', function(){
			expect(element(by.binding('pageHeaderText')).getText()).toBe('Events');
		});
	});

});