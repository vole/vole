define(function(require) {

  var Backbone = require('backbone');

  return Backbone.Model.extend({

    url: '/api/users?is_my_user=true',

    parse: function(response) {
      return response.users[0];
    }

  });

});
