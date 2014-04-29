define(function(require) {

  var $ = require('jquery');
  var Handlebars = require('handlebars');
  var Backbone = require('backbone');
  var format = require('lib/format');

  var HeaderView = require('app/views/header');

  return Backbone.View.extend({

    initialize: function() {
      this.model.on('change:btsync', this.alert.bind(this));
    },

    fetchLayout: function(callback) {
      var layout = format('text!tmpl/layouts/%s.hbs', vole.config.get('ui_layout'));

      require([layout], function(template) {
        callback(Handlebars.compile(template));
      });
    },

    loadTheme: function() {
      var theme = $('<link>', {
        id: 'theme',
        rel: 'stylesheet',
        href: format('/css/themes/%s.css', vole.config.get('ui_theme'))
      });

      $('head').append(theme);
    },

    alert: function() {
      if (this.model.get('btsync')) {
        this.$('.js-alert').show().text('Lost connection to Bittorrent Sync...');
      }
      else {
        this.$('.js-alert').hide();
      }
    },

    render: function() {
      this.fetchLayout(function(layout) {
        this.$el.html(layout());

        new HeaderView().render().$el.appendTo('#header');

        this.loadTheme();
        vole.events.trigger('app.view.rendered');
      }.bind(this));
    }

  });

});
