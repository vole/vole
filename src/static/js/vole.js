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
    user: DS.attr('string')
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
      //controller.set('my_user', App.usersStore.findMyUser());
      controller.set('posts', App.Post.find());
    }
  })

})(jQuery, Ember);
