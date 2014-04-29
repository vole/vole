define(function(require) {

	var $ = require('jquery');

	return {

		posts : function (params) {
			var dfd = $.Deferred();

			$.get('/api/posts', params || {}).done(function (result) {
				if (result && result.posts) {
					dfd.resolve(result.posts);
				}
				else {
					dfd.reject();
				}
			});

			return dfd.promise();
		},

		user : function () {
			var dfd = $.Deferred();

			$.get('/api/users?is_my_user=true').done(function (result) {
				if (result && result.users.length) {
					dfd.resolve(result.users[0]);
				}
				else {
					dfd.reject();
				}
			});

			return dfd.promise();
		},

		createUser : function (name, email) {
			var body = JSON.stringify({
				user : {
					name : name,
					email : email
				}
			});

			return $.ajax({
				type: 'POST',
				url: '/api/users',
				data: body,
				contentType: 'application/json; charset=utf-8',
				dataType: 'json'
			});
		},

		post : function (title) {
			var body = JSON.stringify({
				post : {
					title : title
				}
			});

			return $.ajax({
				type: 'POST',
				url: '/api/posts',
				data: body,
				contentType: 'application/json; charset=utf-8',
				dataType: 'json'
			});
		},

		deletePost : function (id) {
			return $.ajax({
				type: 'DELETE',
				url: '/api/posts/' + id
			});
		},

		addFriend: function(key) {
			return $.ajax({
				type: 'POST',
				url: '/api/friend',
				data: {
					key: key
				}
			});
		}

	};

});
