define(function(require) {

  var Backbone = require('backbone');

  return Backbone.View.extend({

    className: 'profile',

    template: Handlebars.compile(require('text!tmpl/profile.hbs')),

    initialize: function() {
      this.model.on('change', this.render.bind(this));
    },

    render: function() {
      this.$el.html(this.template(this.model.attributes));
      return this;
    }

  });

});
