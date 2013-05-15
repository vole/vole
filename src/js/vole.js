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

  //-------------------------
  // Stores
  //-------------------------
  App.PostsStore = Ember.Object.extend({
    _items: [],

    all: function() {
      return this.get('_items');
    }.property('_items'),

    findAll: function() {
      var self = this;
      var request = $.ajax('/rest/posts.json');

      request.done(function(data, textStatus, jqXHR) {
        if (jqXHR.status === 200) {
          var items = self.get('_items');
          items.slice(0,0);
          data.posts.map(function(post) {
            items.pushObject(App.PostModel.create(post));
          });
        }
      });

      request.fail(function(jqXHR, textStatus, errorThrown) {
        throw new Error('Unable to load posts: ' + textStatus);
      });

      return this.get('all');
    }
  });
  App.postsStore = App.PostsStore.create();

  //-------------------------
  // Views
  //-------------------------

  //-------------------------
  // Controllers
  //-------------------------
  App.ProfileController = Ember.ArrayController.extend({
  })

  //-------------------------
  // Router
  //-------------------------
  App.Router.map(function() {
    this.resource("index", {path: "/"});
    this.resource("profile", {path: "/profile"});
  });

  App.ProfileRoute = Ember.Route.extend({
    model: function() {
      return App.postsStore.findAll();
    }
  })

})(jQuery, Ember);
