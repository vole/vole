define(function(require) {

  var Backbone = require('backbone');

  return Backbone.View.extend({

    template: Handlebars.compile(require('text!tmpl/friend.hbs')),

    events: {
      'click a': 'click'
    },

    click: function(e) {
      e.preventDefault();
      Backbone.history.navigate(this.$('a').attr('href'), true);
    },

    render: function() {
      this.$el.html(this.template(this.model.attributes));
      return this;
    }

  });

});
