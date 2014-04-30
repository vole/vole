define(function(require) {

  var Handlebars = require('handlebars');
  var BaseView = require('app/views/base');
  var PostView = require('app/views/post');

  return BaseView.extend({

    className: 'posts',

    template: Handlebars.compile(require('text!tmpl/posts.hbs')),

    events: {
      'click .js-load-more': 'loadMore'
    },

    initialize: function() {
      this.collection.on('add', this.renderPost, this);
      this.collection.fetch();
    },

    loadMore: function() {
      this.collection.fetch({
        data: {
          before: this.collection.last().get('id')
        }
      });
    },

    renderPost: function(post) {
      this.subView('.posts-list', new PostView({
        tagName: 'li',
        model: post
      }));
    },

    render: function() {
      this.$el.html(this.template());
      return this;
    }

  });

});
