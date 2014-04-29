define(function(require) {

  var Handlebars = require('handlebars');
  var Backbone = require('backbone');

  var PostsView = require('app/views/posts');
  var FriendsView = require('app/views/friends');

  var Posts = require('app/collections/post');
  var Friends = require('app/collections/friend');

  return Backbone.View.extend({

    template: Handlebars.compile(require('text!tmpl/timeline.hbs')),

    events: {
      'click .js-compose': 'compose'
    },

    initialize: function(options) {
      this.options = options;
    },

    render: function() {
      this.$el.html(this.template());

      new PostsView({
        el: this.$('#posts'),
        collection: new Posts([], this.options.user)
      }).render();

      new FriendsView({
        el: this.$('#friends'),
        collection: new Friends()
      }).render();

      return this;
    }

  });

});
