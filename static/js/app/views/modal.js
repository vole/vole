define(function(require) {

  var $ = require('jquery');
  var Backbone = require('backbone');
  var Handlebars = require('handlebars');

  return Backbone.View.extend({

    className: 'modal',

    template: '',

    events: {
      'click .js-close': 'closeButton',
      'click': 'closeBackground'
    },

    initialize: function(options) {
      if (options.template) {
        this.template = options.template;
      }
    },

    ensureOverlay: function() {
      var overlay = $('#modal-overlay');

      if (!overlay.length) {
        overlay = $('<div>', {
          id: 'modal-overlay'
        }).appendTo('body');
      }
    },

    open: function() {
      this.ensureOverlay();

      $('#modal-overlay').show();
      this.$el.show();
    },

    closeButton: function() {
      $('#modal-overlay').hide();
      this.remove();
    },

    closeBackground: function(e) {
      if (e.target !== e.currentTarget) {
        return;
      }

      $('#modal-overlay').hide();
      this.remove();
    },

    render: function() {
      var template = Handlebars.compile(this.template);

      var content = $('<div>', {
        'class': 'modal-content',
        html: template()
      });

      this.$el.html(content).appendTo('body');

      return this;
    }

  });

});
