define(function(require) {

  var Backbone = require('backbone');

  return Backbone.Router.extend({

    routes: {
      //'': require('app/routes/timeline'),
      'timeline': require('app/routes/timeline'),
      'timeline/:action': require('app/routes/timeline'),
      'timeline/friend/:friend': require('app/routes/timeline'),

      'compose': require('app/routes/compose'),
      'compose/:id': require('app/routes/compose'),

      'profile': require('app/routes/profile'),
      'preferences': require('app/routes/preferences')
    }

  });

});
