define(function(require) {

  var _ = require('underscore');
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
      this.collection.on('sync', this.renderResults.bind(this));
      this.collection.fetch();
    },

    search: _.debounce(function() {
      this.collection.fetch({
        data: {
          query: this.$('#query').val()
        }
      });
    }, 100),

    renderResults: function() {
      this.$('ul').empty();

      this.collection.each(function(user) {
        this.subView('ul', new FriendView({
          model: user,
          tagName: 'li'
        }));
      }.bind(this));
    },

    render: function() {
      this.$el.html(this.template);
      return this;
    }

  });

});
