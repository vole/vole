define(function(require) {

  var BaseView = require('app/views/base');
  var Handlebars = require('handlebars');

  return BaseView.extend({

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
