define(function(require) {

  var Dropdown = require('app/views/dropdown');
  var api = require('app/api');

  return Dropdown.extend({

    events: {
      'click a': 'add',
      'click': 'click'
    },

    click: function(e) {
      e.stopPropagation();
    },

    add: function(e) {
      e.preventDefault();

      var key = this.$('input').val();

      // TODO: error handling!
      api.addFriend(key).always(this.close.bind(this));

      Backbone.history.navigate('/timeline', true);

      this.$('input').val('');
    }

  });

});
