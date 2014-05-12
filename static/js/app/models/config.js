define(function(require) {

  var Backbone = require('backbone');

  return Backbone.Model.extend({

    url: '/config',

    defaults: {
      ui_logging: 'info'
    }

  });

});
