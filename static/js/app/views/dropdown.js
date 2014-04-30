define(function(require) {

  var _ = require('underscore');
  var Backbone = require('backbone');
  var Handlebars = require('handlebars');

  return Backbone.View.extend({

    className: 'dropdown',

    events: {
      'click': 'click'
    },

    initialize: function(options) {
      options = options || {};
      this.position = options.position || 'left';
      this.target = options.target;
      this.template = options.template || '';

      $(document).on('click', this.close.bind(this));
    },

    click: function(e) {
      e.stopPropagation();
    },

    open: function() {
      var h = this.target.height();
      var w = this.target.width();
      var offset = this.target.offset();

      var top = 0;
      var left = 0;

      switch (this.position) {
        case 'top':
        case 'bottom':
        case 'right':
        case 'left':
        default:
          top = offset.top;
          left = offset.left - this.width();
          break;
      }

      // Set the dropdown's position.
      this.$el.css({
        top: top + 'px',
        left: left + 'px'
      });

      // Reveal the dropdown.
      this.$el.show();
    },

    close: function() {
      this.$el.hide();
    },

    toggle: function() {
      this.$el.toggle();
    },

    width: function() {
      return this.$el.width();
    },

    height: function() {
      return this.$el.height();
    },

    render: function() {
      var template = Handlebars.compile(this.template);

      var content = $('<div>', {
        'class': 'dropdown-content',
        html: template()
      });

      this.$el.html(content).appendTo('body');

      return this;
    }

  });

});
