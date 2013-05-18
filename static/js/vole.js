(function ($, Ember) {
  var cl = console.log.bind(console);
  var App = Ember.Application.create({
    LOG_TRANSITIONS: true,
    rootElement: '#ember-container'
  });
  window.App = App;

  //-------------------------
  // Store
  //-------------------------
  App.Store = DS.Store.extend({
    revision: 12,
    adapter: DS.RESTAdapter
  });

  DS.RESTAdapter.reopen({
    namespace: 'api'
  });

  //-------------------------
  // Models
  //-------------------------
  App.Post = DS.Model.extend({
    title: DS.attr('string'),
    user: DS.attr('string'),
    created: DS.attr('string'),
    created_pretty: DS.attr('string')
  });

  App.User = DS.Model.extend({
    key: DS.attr('string'),
    hash: DS.attr('string'),
    user: DS.attr('string'),
    display_name: DS.attr('string'),
    is_my_user: DS.attr('boolean')
  });

  //-------------------------
  // Views
  //-------------------------
  App.PostsView = Ember.View.extend({
    templateName: 'posts'
  });

  //-------------------------
  // Controllers
  //-------------------------
  App.ApplicationController = Ember.Controller.extend({
    needs: ['posts', 'users']
  });

  App.IndexController = Ember.Controller.extend({
    needs: ['posts', 'users'],
    new_post: '',

    createNewPost: function() {
      var self = this;
      var my_user = this.get('controllers.users.myUser.firstObject.user');

      var newpost = App.Post.createRecord({
        user: my_user,
        title: this.get('new_post')
      });
      newpost.on('didCreate', function() {
        self.set('new_post', '');
      });
      newpost.get('transaction').commit();
    }
  });

  App.ProfileController = Ember.Controller.extend({
    needs: ['posts', 'users'],
  });

  App.UsersController = Ember.ArrayController.extend({
    // This is set to a FilteredRecordArray by the router. Just use the
    // first object in the array.
    myUser: null
  });

  App.PostsController = Ember.ArrayController.extend({
    filterByUser: [],

    filteredPosts: function() {
      if (this.get('filterByUser.length') > 0) {
        var filterUser = this.get('filterByUser.firstObject.user');
        if (filterUser) {
          return this.get('content').filterProperty('user', filterUser);
        }
      }
      else {
        return this.get('content');
      }
    }.property('content.[]', 'filterByUser.[]')
  });

  //-------------------------
  // Router
  //-------------------------
  App.Router.map(function() {
    this.resource('index', {path: '/'});
    this.resource('profile', {path: '/profile'});
  });

  App.ApplicationRoute = Ember.Route.extend({
    setupController: function(controller) {
      controller.set('controllers.posts.content', App.Post.find());
      controller.set('controllers.users.content', App.User.find());
      controller.set('controllers.users.myUser', App.User.filter(function(user) {
        if (user.get('is_my_user')) {
          return true;
        }
      }));
    }
  });

  App.IndexRoute = Ember.Route.extend({
    setupController: function(controller) {
      var postsController = controller.get('controllers.posts');
      postsController.set('filterByUser.[]', []);
    }
  });

  App.ProfileRoute = Ember.Route.extend({
    setupController: function(controller) {
      var postsController = controller.get('controllers.posts');
      var usersController = controller.get('controllers.users');
      postsController.set('filterByUser.[]', usersController.get('myUser.firstObject'));
    },
  });


})(jQuery, Ember);
