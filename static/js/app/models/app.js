define(function(require) {

  var Backbone = require('backbone');

  return Backbone.Model.extend({

    defaults: {
      btsync: true
    },

    url: '/status',

    initialize: function() {
      this.interval = setInterval(this.fetch.bind(this), 5000);
    }

  });

});
