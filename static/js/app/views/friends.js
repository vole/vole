define(function(require) {

  var _ = require('underscore');
  var Backbone = require('backbone');
  var Handlebars = require('handlebars');

  var FriendView = require('app/views/friend');
  var Friends = require('app/collections/friend');

  var api = require('app/api');

  return Backbone.View.extend({

    template: Handlebars.compile(require('text!tmpl/friends.hbs')),

    events: {
      'keyup .js-search': 'search',
      'click .js-add': 'toggleAdd',
      'click .js-add-button': 'add'
    },

    initialize: function() {
      this.collection = new Friends();
      this.collection.on('sync', this.renderResults.bind(this));
      //this.updateInterval = setInterval(this.search.bind(this), 5000);
      this.search();
    },

    parse: function(response) {
      return response.users;
    },

    search: _.debounce(function() {
      this.collection.fetch({
        data: {
          query: this.$('#query').val()
        }
      });
    }, 100),

    toggleAdd: function() {
      this.$('form').toggle();
    },

    add: function(e) {
      e.preventDefault();

      var key = this.$('#add-friend').val();

      api.addFriend(key).done(this.search.bind(this));

      this.toggleAdd();

      this.$('#add-friend').val('');
    },

    renderResults: function() {
      this.$('ul').empty();

      this.collection.each(function(user) {
        this.$('ul').append(
          new FriendView({ model: user, tagName: 'li' }).render().el
        );
      }.bind(this));
    },

    render: function() {
      this.$el.html(this.template);
      return this;
    }

  });

});
