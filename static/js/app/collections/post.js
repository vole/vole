define(function(require) {

  var Backbone = require('backbone');
  var Post = require('app/models/post');

  return Backbone.Collection.extend({

    url: function() {
      var _url = '/api/posts';
      if (this.user) {
        _url += '?user=' + this.user;
      }
      return _url;
    },

    initialize: function(models, options) {
      options = options || {};
      this.user = options.user;
    },

    parse: function(response) {
      return response.posts;
    },

    model: Post

  });

});
