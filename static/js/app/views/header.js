define(function(require) {

  var Backbone = require('backbone');

  var Modal = require('app/views/modal');
  var Add = require('app/views/add');

  return Backbone.View.extend({

    className: 'header',

    template: Handlebars.compile(require('text!tmpl/header.hbs')),

    events: {
      'click a[href!=#]': 'navigate',
      'click .js-info': 'info',
      'click .js-add': 'add'
    },

    navigate: function(e) {
      e.preventDefault();
      var href = $(e.target).attr('href');
      Backbone.history.navigate(href, true);
    },

    add: function(e) {
      e.preventDefault();
      e.stopPropagation();

      if (this.dropdown) {
        return this.dropdown.toggle();
      }

      this.dropdown = new Add({
        target: this.$('.js-add'),
        template: require('text!tmpl/dropdowns/add.hbs')
      });

      this.dropdown.render().open();
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
