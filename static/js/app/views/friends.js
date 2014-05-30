define(function(require) {

  var Handlebars = require('handlebars');
  var BaseView = require('app/views/base');
  var FriendView = require('app/views/friend');

  return BaseView.extend({

    className: 'friends',

    template: Handlebars.compile(require('text!tmpl/friends.hbs')),

    events: {
      'keyup .js-search': 'search'
    },

    initialize: function() {
      this.collection.on('sync', this.render, this);
      this.collection.fetch();
    },

    search: function() {
      var query = this.$('#query').val().toLowerCase();

      this.subViews().forEach(function(view) {
        var name = view.model.get('name').toLowerCase();
        view.$el.toggle(name.indexOf(query) > -1);
      });
    },

    render: function() {
      this.$el.html(this.template);

      this.collection.forEach(function(user) {
        var view = new FriendView({
          model: user,
          tagName: 'li'
        });

        this.subView('ul', view);
      }, this);

      return this;
    }

  });

});
