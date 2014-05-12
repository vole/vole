define(function(require) {

  var Backbone = require('backbone');

  return Backbone.Model.extend({

    defaults: {
      step: 1,
      name: null,
      key: null,
      avatar: null
    }

  });

});
