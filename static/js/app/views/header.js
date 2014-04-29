define(function(require) {

  var Backbone = require('backbone');

  var Modal = require('app/views/modal');

  return Backbone.View.extend({

    className: 'header',

    template: Handlebars.compile(require('text!tmpl/header.hbs')),

    events: {
      'click a[href!=#]': 'navigate',
      'click .js-info': 'info'
    },

    navigate: function(e) {
      e.preventDefault();
      var href = $(e.target).attr('href');
      Backbone.history.navigate(href, true);
    },

    info: function(e) {
      e.preventDefault();

      var modal = new Modal({
        template: require('text!tmpl/modals/info.hbs')
      });

      modal.render().open();
    },

    render: function() {
      this.$el.html(this.template());
      return this;
    }

  });

});
