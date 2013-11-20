
define([
	'flight/component',
	'jquery',
	'app/events',
	'app/api'
], function (component, $, events, API) {

	function post () {

		this.defaultAttrs({
			deletePostSelector: '.js-delete'
		});

		this.deletePost = function () {
			if (confirm('Are you sure you want to delete this post?')) {
				API.deletePost(this.id());
				this.$node.fadeOut();
				events.trigger(events.POST_DELETED, this.id());
			}
		};

		this.id = function () {
			return this.$node.data('id');
		};

		this.after('initialize', function () {
			this.on('click', {
				deletePostSelector: this.deletePost
			});
		});

	}

	return component(post);

});
