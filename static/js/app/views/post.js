define(function(require) {

  var Backbone = require('backbone');
  var Handlebars = require('handlebars');

  return Backbone.View.extend({

    className: 'post',

    events: {
      'click .js-delete': 'delete'
    },

    template: Handlebars.compile(require('text!tmpl/post.hbs')),

    initialize: function() {
      this.model.on('destroy', this.fadeOut.bind(this));
    },

    delete: function() {
      if (confirm('You sure?')) {
        this.model.destroy();
      }
    },

    fadeOut: function() {
      this.$el.fadeOut(this.remove.bind(this));
    },

    render: function() {
      this.$el.html(this.template(this.model.attributes));
      return this;
    }

  });

});
