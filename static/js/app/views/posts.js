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
      this.collection.on('sync', this.onSync, this);
      this.collection.fetch();
    },

    onSync: function(collection, models) {
      if (models.length < vole.config.get('ui_pageSize')) {
        this.disableLoadMore()
      }
    },

    loadMore: function() {
      this.collection.fetch({
        remove: false,
        data: {
          before: this.collection.last().get('id')
        }
      });
    },

    disableLoadMore: function() {
      this.$('.js-load-more').addClass('disabled');
    },

    renderPost: function(post) {
      this.$('.js-load-more').show();
      this.$('.empty').hide();

      this.subView('.js-posts', new PostView({
        tagName: 'li',
        model: post
      }));
    },

    renderPosts: function() {
      this.collection.each(this.renderPost, this);
    },

    render: function() {
      this.$el.html(this.template());
      this.renderPosts();
      return this;
    }

  });

});
