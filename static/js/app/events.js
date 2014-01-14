
define(['jquery'], function ($) {

	var $doc = $(document);

	return {
		// View the user's personal timeline.
		VIEW_MY_PROFILE: 'VIEW_MY_PROFILE',

		// View the user's home timeline.
		VIEW_HOME: 'VIEW_HOME',

		// The user tried to view their profile, but they do not have one.
		PROFILE_NOT_FOUND: 'PROFILE_NOT_FOUND',

		// A new profile has been created.
		PROFILE_CREATED: 'PROFILE_CREATED',

		POST_DELETED: 'POST_DELETED',

		WRITE_POST: 'WRITE_POST',

		trigger: function () {
			var args = Array.prototype.slice.call(arguments);
			$doc.trigger.apply($doc, args);
		},

		on: function () {
			var args = Array.prototype.slice.call(arguments);
			$doc.on.apply($doc, args);
		}
	};

});
