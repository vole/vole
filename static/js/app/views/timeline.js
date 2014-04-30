define(function(require) {

  var Handlebars = require('handlebars');
  var Backbone = require('backbone');

  var BaseView = require('app/views/base');
  var PostsView = require('app/views/posts');
  var FriendsView = require('app/views/friends');

  var Posts = require('app/collections/post');
  var Friends = require('app/collections/friend');

  return BaseView.extend({

    template: Handlebars.compile(require('text!tmpl/timeline.hbs')),

    events: {
      'click .js-compose': 'compose'
    },

    initialize: function(options) {
      this.options = options;
    },

    render: function() {
      this.$el.html(this.template());

      this.subView('#posts', new PostsView({
        collection: new Posts([], this.options)
      }));

      this.subView('#friends', new FriendsView({
        collection: new Friends()
      }));

      return this;
    }

  });

});
