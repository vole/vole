
define(['jquery'], function ($) {

	var API = {

		init : function () {
			this.getToken(function () {

			});

			return this;
		},

		baseUrl : function () {
			return 'http://127.0.0.1:26085/gui';
		},

		request : function (path, params, callback) {
			var url = this.baseUrl() + '/' + path + '?' + $.param(params);
			return $.get(url, callback || $.noop);
		},

		getToken : function (callback) {
			this.request('token.html', { t : Date.now() }).done(function (data) {
				callback($(data).find('#token'));
			});
		}

	};

	return API;

});
