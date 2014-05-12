define(function(require) {

  var Handlebars = require('handlebars');
  var Backbone = require('backbone');

  return Backbone.View.extend({

    className: 'error',

    template: Handlebars.compile(require('text!tmpl/error.hbs')),

    render: function() {
      this.$el.html(this.template());
      return this;
    }

  });

});
