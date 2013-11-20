
define([
	'flight/component',
	'app/api',
	'app/events',
	'lib/dropzone'
], function (component, API, events, Dropzone) {

	function postbox () {

		this.defaultAttrs({
			textareaSelector : 'textarea',
			buttonSelector : '.modal-footer button',
			filesSelector: '#post-files',
			dropzoneSelector: '.dropzone',
			uploadPath: '/file/upload'
		});

		this.post = function () {
			var files = this.files.map(function (file) {
				return {
					name: file.name,
					hash: file.hash,
					type: file.type,
					size: file.size
				};
			});

			API
				.post(this.text(), files)
				.done(this.close.bind(this));
		};

		this.keyup = function () {
			this.select('buttonSelector')
				.prop('disabled', this.text().length === 0);
		};

		this.text = function () {
			return this.select('textareaSelector').val();
		};

		this.open = function () {
			this.$node.modal('show');
		};

		this.close = function () {
			this.$node.modal('hide');
		};

		this.reset = function () {
			this.select('textareaSelector').val('');
			this.files = [];
			this.dropzone.removeAllFiles();
		};

		this.fileUploaded = function (file, response) {
			// Augment the file object with its hash.
			file.hash = response.hash;
			this.files.push(file);
		};

		this.after('initialize', function () {
			this.on('click', {
				buttonSelector : this.post
			});

			this.on('keyup', {
				textareaSelector : this.keyup
			});

			events.on(events.WRITE_POST, this.open.bind(this));

			Dropzone.autoDiscover = false;

			this.files = [];

			this.dropzone = new Dropzone(this.select('dropzoneSelector').get(0), {
				url: this.attr.uploadPath,
				addRemoveLinks: true,
				dictDefaultMessage: 'Click or drop to attach files'
			});

			this.dropzone.on('success', this.fileUploaded.bind(this));
		});

		this.after('close', this.reset);
	}

	return component(postbox);

});
