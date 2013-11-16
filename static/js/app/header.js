
define([
	'flight/component',
	'app/api',
	'app/events'
], function (component, API, events) {

	function header () {

		this.defaultAttrs({
			username : '#my-username-title'
		});

		this.renderUsername = function () {
			API.user().done(function (user) {
				this.select('username').text(user.name);
			}.bind(this));
		};

		this.after('initialize', function () {
			this.renderUsername();

			$(document).on(events.PROFILE_CREATED, this.renderUsername.bind(this));
		});

	}

	return component(header);

});
