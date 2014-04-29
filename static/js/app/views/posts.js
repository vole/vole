define(function(require) {

  var Backbone = require('backbone');
  var Handlebars = require('handlebars');
  var PostView = require('app/views/post');
  var Posts = require('app/collections/post');
  var Spinner = require('lib/spin');

  return Backbone.View.extend({

    className: 'posts',

    template: Handlebars.compile(require('text!tmpl/posts.hbs')),

    events: {
      'click .js-load-more': 'loadMore'
    },

    initialize: function() {
      this.collection.on('sync', this.render.bind(this));
      this.collection.fetch();
      this.interval = setInterval(this.checkForUpdates.bind(this), 5000);
    },

    checkForUpdates: function() {
      var posts = new Posts();

      posts.on('sync', function() {
        var index = posts.indexOf(this.collection.first());

        if (index > 0) {
          var newPosts = posts.slice(0, index - 1);
          console.log('new posts available:', newPosts);
        }
      }.bind(this));

      posts.fetch();
    },

    loadMore: function(e) {
      e.preventDefault();

      this.renderSpinner();
      this.$('.js-load-more').addClass('disabled');

      var posts = new Posts();

      posts.on('sync', function() {
        this.collection.add(posts.models);
        this.renderPosts(posts);
        this.$('.js-load-more').removeClass('disabled');
      }.bind(this));

      posts.fetch({
        data: {
          before: this.collection.last().get('id')
        }
      });
    },

    renderPosts: function(collection) {
      var loadMoreButton = this.$('.load-more');

      this.$('.spinner').remove();

      collection.each(function(post) {
        var view = new PostView({
          model: post
        });

        view.render().$el.insertBefore(loadMoreButton);
      }.bind(this));
    },

    renderSpinner: function() {
      var spinner = new Spinner({
        width: 1,
        height: 1,
        radius: 5
      }).spin();

      this.$('.js-load-more').before(spinner.el);
    },

    render: function() {
      this.$el.html(this.template());
      this.renderSpinner();
      this.renderPosts(this.collection);
      return this;
    }

  });

});
