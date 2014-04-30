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
      this.model.on('sync', this.alert.bind(this));
    },

    loadTheme: function() {
      var href = format('/css/themes/%s.css', vole.config.get('ui_theme'));
      var theme = $('<link id="theme" rel="stylesheet">').attr('href', href);
      $('head').append(theme);
    },

    // TODO: Make this its own view.
    alert: function() {
      if (!this.model.get('btsync')) {
        this.$('.js-alert').slideDown();
      }
      else {
        this.$('.js-alert').slideUp();
      }
    },

    render: function() {
      this.$el.html(this.template(this.model.attributes));

      // Render the header. It's a persistant view, always visible.
      var header = new HeaderView().render();
      this.$('#header').append(header.el);

      // Load the user's configured theme.
      this.loadTheme();

      vole.events.trigger('app.view.rendered');
    },

    setContentView: function(contentView) {
      if (this.contentView && this.contentView._subViews) {
        var toKill = this.contentView._subViews || [];

        while (toKill.length) {
          var view = toKill.shift();
          view.trigger('kill');

          if (view._subViews) {
            toKill = toKill.concat(view._subViews || []);
          }
        }
      }

      this.contentView = contentView;

      this.$('#content').html(contentView.render().el);
    }

  });

});
