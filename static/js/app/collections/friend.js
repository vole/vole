define(function(require) {

  var Backbone = require('backbone');
  var Friend = require('app/models/friend');

  return Backbone.Collection.extend({

    url: '/api/users',

    model: Friend

  });

});
