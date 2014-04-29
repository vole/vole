define(function(require) {

  var Posts = require('app/collections/post');

  return Posts.extend({

    url: function() {
      return '/api/drafts';
    }

  });

});
