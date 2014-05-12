define(function(require) {

  var Backbone = require('backbone');

  return Backbone.Router.extend({

    initialize: function() {
      this.route('install', 'install', require('app/routes/install'));

      this.route('timeline', 'timeline', require('app/routes/timeline'));
      this.route('timeline/:friend', 'friend\'s timeline', require('app/routes/timeline'));

      this.route('compose', 'compose', require('app/routes/compose'));
      this.route('compose/:id', 'edit', require('app/routes/compose'));

      this.route('profile', 'profile', require('app/routes/profile'));
      this.route('preferences', 'preferences', require('app/routes/preferences'));
      this.route('error', 'error', require('app/routes/error'));
    }

  });

});
