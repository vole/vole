
define([
	'flight/component',
	'jquery',
	'app/events',
	'app/api'
], function (component, $, events, API) {

	function createProfile () {

		this.defaultAttrs({
			buttonSelector: 'button',
			nameSelector: '#create-profile-name',
			emailSelector: '#create-profile-email'
		});

		this.createProfile = function () {
			this.$node.modal('show');
		};

		this.create = function () {
			var email = this.select('emailSelector').val();
			var name = this.select('nameSelector').val();

			API.createUser(name, email).done(function () {
				this.$node.modal('hide');
				$(document).trigger(events.PROFILE_CREATED);
				$(document).trigger(events.VIEW_MY_PROFILE);
			}.bind(this));
		};

		this.keyup = function () {
			var disabled = (this.select('nameSelector').val().length === 0);
			this.select('buttonSelector').prop('disabled', disabled);
		};

		this.after('initialize', function () {
			$(document).on(events.PROFILE_NOT_FOUND, $.proxy(this.createProfile, this));

			this.$node.modal({
				show: false
			});

			this.on('click', {
				buttonSelector: this.create
			});

			this.on('keyup', {
				nameSelector : this.keyup
			});
		});
	}

	return component(createProfile);

});
