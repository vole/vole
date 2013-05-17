(function ($, Ember) {
  var cl = console.log.bind(console);
  var App = Ember.Application.create({
    LOG_TRANSITIONS: true,
    rootElement: '#ember-container'
  });
  //Ember.LOG_BINDINGS = true;
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
    created: DS.attr('string')
  });

  App.User = DS.Model.extend({
    user: DS.attr('string'),
    display_name: DS.attr('string')
  });

  //-------------------------
  // Views
  //-------------------------

  //-------------------------
  // Controllers
  //-------------------------
  App.ProfileController = Ember.ObjectController.extend({
    posts: [],
    my_user: [],
    my_user_name: '',

    myPosts: function() {
      var my_user = this.get('my_user');
      if (my_user.get('length') > 0) {
        this.set('my_user_name', my_user.objectAt(0).get('user'));
      }
      return this.get('posts').filterProperty('user', this.get('my_user_name'));
    }.property('my_user.@each.user', 'posts.@each'),
  });

  App.IndexController = Ember.ObjectController.extend({
    posts: [],
    my_user: [],
    my_user_name: '',
    new_post: '',

    createNewPost: function() {
      var my_user = this.get('my_user');
      if (my_user.get('length') > 0) {
        this.set('my_user_name', my_user.objectAt(0).get('user'));
      }

      var newpost = App.Post.createRecord({
        user: this.get('my_user_name'),
        title: this.get('new_post')
      });
      newpost.get('transaction').commit();
    }
  });

  //-------------------------
  // Router
  //-------------------------
  App.Router.map(function() {
    this.resource('index', {path: '/'});
    this.resource('profile', {path: '/profile'});
  });

  App.ProfileRoute = Ember.Route.extend({
    setupController: function(controller) {
      controller.set('my_user', App.User.find());
      controller.set('posts', App.Post.find());
    }
  });

  App.IndexRoute = Ember.Route.extend({
    setupController: function(controller) {
      controller.set('my_user', App.User.find());
      controller.set('posts', App.Post.find());
    }
  })

})(jQuery, Ember);
