define(function(require) {

  var Dropdown = require('app/views/dropdown');
  var Friend = require('app/models/friend');

  return Dropdown.extend({

    events: {
      'click .js-add': 'add',
      'click .js-close': 'close',
      'click': 'click'
    },

    click: function(e) {
      e.stopPropagation();
    },

    add: function(e) {
      e.preventDefault();

      var friend = new Friend({
        id: this.$('input').val()
      });

      this.clear();

      this.showLoading();

      friend.save({}, {
        success: this.success.bind(this),
        error: this.error.bind(this)
      });
    },

    clear: function() {
      this.$('input').val('').removeClass('error');
      this.$('.js-error').hide().text('');
      this.$('.js-spin').empty();
    },

    success: function() {
      this.hideLoading();
      this.remove();
      Backbone.history.navigate('/timeline', true);
    },

    error: function() {
      this.hideLoading();
      this.$('input').addClass('error');
      this.$('.js-error').show().text('Error adding user.');
      this.$('.js-spin').empty();
    }

  });

});
