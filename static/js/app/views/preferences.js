define(function(require) {

  var Handlebars = require('handlebars');
  var Backbone = require('backbone');

  return Backbone.View.extend({

    className: 'preferences',

    template: Handlebars.compile(require('text!tmpl/preferences.hbs')),

    events: {
      'change .js-theme': 'theme'
    },

    theme: function() {
      var theme = this.$('.js-theme').val();
      $('#theme').attr('href', '/css/themes/' + theme + '.css');
    },

    render: function() {
      this.$el.html(this.template(vole.config.attributes));
      return this;
    }

  });

});
