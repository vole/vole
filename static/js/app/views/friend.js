define(function(require) {

  var Backbone = require('backbone');
  var BaseView = require('app/views/base');

  return BaseView.extend({

    className: 'friend',

    template: Handlebars.compile(require('text!tmpl/friend.hbs')),

    events: {
      'click a': 'click'
    },

    // TODO: Instead of navigating, just re-render the posts view.
    // Maybe using events somehow?
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
