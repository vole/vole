define(function(require) {

  var format = require('lib/format');
  var md5 = require('lib/md5');
  var Backbone = require('backbone');

  return Backbone.Model.extend({

    url: '/api/users?is_my_user=true',

    defaults: {
      email: ''
    },

    initialize: function() {
      this.on('change:email', this.gravatar, this);
    },

    gravatar: function() {
      var hash = md5(this.get('email').trim().toLowerCase());
      var image = format('http://www.gravatar.com/avatar/%s', hash);
      this.set('avatar', image);
    }

  });

});
