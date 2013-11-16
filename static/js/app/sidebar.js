
define([
	'flight/component',
	'jquery',
	'app/events'
], function (component, $, events) {

	function sidebar () {

		this.defaultAttrs({
			viewProfileSelector: '.view-my-profile',
			viewHomeSelector: '.view-home',
			menuItemSelector: 'li'
		});

		this.viewProfile = function () {
			$(document).trigger(events.VIEW_MY_PROFILE);
		};

		this.viewHome = function () {
			$(document).trigger(events.VIEW_HOME);
		};

		this.menuItemClicked = function (e, data) {
			this.select('menuItemSelector').removeClass('active');
			$(data.el).addClass('active');
		};

		this.after('initialize', function () {
			this.on('click', {
				viewProfileSelector: this.viewProfile,
				viewHomeSelector: this.viewHome,
				menuItemSelector: this.menuItemClicked
			});
		});

	}

	return component(sidebar);

});
