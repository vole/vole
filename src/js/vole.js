(function ($, Ember) {
  var cl = console.log.bind(console);
  var App = Ember.Application.create({
    LOG_TRANSITIONS: true,
    rootElement: '#ember-container'
  });
  //Ember.LOG_BINDINGS = true;
  window.App = App;

  //-------------------------
  // Models
  //-------------------------
  App.PostModel = Ember.Object.extend({
    title: '',
    date: '',
    text: '',
    user: ''
  });

  App.UserModel = Ember.Object.extend({
    user: '',
    display_name: ''
  });

  //-------------------------
  // Stores
  //-------------------------
  App.PostsStore = Ember.Object.extend({
    posts: [],

    all: function() {
      return this.get('posts');
    }.property('posts'),

    findAll: function() {
      var self = this;
      var request = $.ajax('/api/posts');

      request.done(function(data, textStatus, jqXHR) {
        if (jqXHR.status === 200) {
          var items = self.get('posts');
          items.slice(0,0);
          data.posts.map(function(post) {
            items.pushObject(App.PostModel.create(post));
          });
        }
      });
      request.fail(function(jqXHR, textStatus, errorThrown) {
        throw new Error('Unable to load posts: ' + textStatus);
      });
      return this.get('posts');
    }
  });
  App.postsStore = App.PostsStore.create();

  App.UsersStore = Ember.Object.extend({
    users: [],
    my_user: [],

    findMyUser: function() {
      var self = this;
      var request = $.ajax('/api/my_user');

      request.done(function(data, textStatus, jqXHR) {
        if (jqXHR.status === 200) {
          self.get('my_user').pushObject(App.UserModel.create(data.users[0]));
        }
      });
      request.fail(function(jqXHR, textStatus, errorThrown) {
        throw new Error('Unable to load my user: ' + textStatus);
      });
      return this.get('my_user');
    }
  });
  App.usersStore = App.UsersStore.create();

  //-------------------------
  // Views
  //-------------------------

  //-------------------------
  // Controllers
  //-------------------------
  App.ProfileController = Ember.ObjectController.extend({
    posts: [],
    my_user: [],

    myPosts: function() {
      var my_user = this.get('my_user');
      var my_user_name = '';
      if (my_user.length > 0) {
        my_user_name = my_user.objectAt(0).get('user');
      }

      return this.get('posts').filterProperty('user', my_user_name);
    }.property('my_user.@each.user', 'posts.@each'),

  })

  //-------------------------
  // Router
  //-------------------------
  App.Router.map(function() {
    this.resource("index", {path: "/"});
    this.resource("profile", {path: "/profile"});
  });

  App.ProfileRoute = Ember.Route.extend({
    setupController: function(controller) {
      controller.set('my_user', App.usersStore.findMyUser());
      controller.set('posts', App.postsStore.findAll());
    }
  })

})(jQuery, Ember);
