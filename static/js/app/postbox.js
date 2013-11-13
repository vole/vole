
define([
	'flight/component',
	'app/api'
], function (component, API) {

	function postbox () {

		this.defaultAttrs({
			textareaSelector : 'textarea',
			buttonSelector : '#post-button'
		});

		this.post = function () {
			var textarea = this.select('textareaSelector');

			API.post(textarea.val()).done(function () {
				textarea.val('');
			}.bind(this));
		};

		this.keyup = function () {
			var disabled = (this.select('textareaSelector').val().length === 0);
			this.select('buttonSelector').prop('disabled', disabled);
		};

		this.after('initialize', function () {
			this.on('click', {
				buttonSelector : this.post
			});

			this.on('keyup', {
				textareaSelector : this.keyup
			});
		});

	}

	return component(postbox);

});
