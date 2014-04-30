define(function(require) {

  var $ = require('jquery');
  var _ = require('underscore');
  var Handlebars = require('handlebars');
  var Backbone = require('backbone');
  var format = require('lib/format');

  var HeaderView = require('app/views/header');

  return Backbone.View.extend({

    template: Handlebars.compile(require('text!tmpl/layouts/default.hbs')),

    initialize: function() {
      this.model.on('change:btsync', this.alert.bind(this));
    },

    loadTheme: function() {
      var theme = $('<link>', {
        id: 'theme',
        rel: 'stylesheet',
        href: format('/css/themes/%s.css', vole.config.get('ui_theme'))
      });

      $('head').append(theme);
    },

    // TODO: Make this its own view.
    alert: function() {
      if (this.model.get('btsync')) {
        this.$('.js-alert').show().text('Lost connection to Bittorrent Sync...');
      }
      else {
        this.$('.js-alert').hide();
      }
    },

    render: function() {
      this.$el.html(this.template());

      new HeaderView().render().$el.appendTo('#header');

      this.loadTheme();

      vole.events.trigger('app.view.rendered');
    },

    setContentView: function(contentView) {
      if (this.contentView && this.contentView._subViews) {
        var toKill = _.values(this.contentView._subViews);

        while (toKill.length) {
          var view = toKill.shift();
          view.trigger('kill');

          if (view._subViews) {
            toKill = toKill.concat(_.values(view._subViews));
          }
        }
      }

      this.contentView = contentView;

      this.$('#content').html(contentView.render().el);
    }

  });

});
