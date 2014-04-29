define(function(require) {

  var Backbone = require('backbone');
  var Friend = require('app/models/friend');

  return Backbone.Collection.extend({

    url: '/api/users',

    parse: function(response) {
      return response.users;
    },

    model: Friend

  });

});
