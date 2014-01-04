
define([
	'flight/component',
	'app/api',
	'lib/marked'
], function (component, API, marked) {

	function postbox () {

		this.defaultAttrs({
			textareaSelector : 'textarea',
			previewSelector  : '#preview-text',
			buttonSelector   : '#post-button'
		});

		this.post = function () {
			var textarea = this.select('textareaSelector');
			this.select('previewSelector').html('');
			API.post(textarea.val()).done(function () {
				textarea.val('');
			}.bind(this));
		};

		this.keyup = function () {
			var textarea = this.select('textareaSelector');
			var disabled = (textarea.val().length === 0);
			var height = textarea.innerHeight();
			var scrollHeight = textarea[0].scrollHeight;

			//Respond to overflows in textarea content
			if( scrollHeight > height ){
				textarea.height(scrollHeight);
			}

			this.select('buttonSelector').prop('disabled', disabled);
			this.select('previewSelector')
				.html(marked(textarea.val()))
				.toggleClass('invisible', disabled);
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
