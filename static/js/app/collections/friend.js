define(function(require) {

  var Backbone = require('backbone');

  return Backbone.Collection.extend({

    url: '/api/users'

  });

});
