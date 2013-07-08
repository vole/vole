
define([
  'app/config',
  'ember',
  'ember-data',
  'lib/marked',
  'plugins/text!app/templates/application.hbs',
  'plugins/text!app/templates/index.hbs',
  'plugins/text!app/templates/posts.hbs',
  'plugins/text!app/templates/profile.hbs',
  'app/templates/helpers',
  'plugins/moment',
  'plugins/resize'
],
function (Config, Ember, DS, marked, applicationTemplate, indexTemplate, postsTemplate, profileTemplate) {
  Ember.TEMPLATES['application'] = Ember.Handlebars.compile(applicationTemplate);
  Ember.TEMPLATES['index'] = Ember.Handlebars.compile(indexTemplate);
  Ember.TEMPLATES['profile'] = Ember.Handlebars.compile(profileTemplate);
  Ember.TEMPLATES['posts'] = Ember.Handlebars.compile(postsTemplate);

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
    created: DS.attr('number'),
    userId: DS.attr('string'),
    userName: DS.attr('string'),
    userAvatar: DS.attr('string'),
    isMyPost: DS.attr('boolean'),
    isRepost: DS.attr('boolean'),
    repostUserName: DS.attr('string'),
    repostUserAvatar: DS.attr('string'),
    repostPostId: DS.attr('string')
  });

  App.User = DS.Model.extend({
    name: DS.attr('string'),
    avatar: DS.attr('string'),
    isMyUser: DS.attr('boolean'),
    email: DS.attr('string')
  });

  //-------------------------
  // Views
  //-------------------------
  App.PostsView = Ember.View.extend({
    templateName: 'posts'
  });

  App.IndexView = Ember.View.extend({
    keyPress: function(event) {
      if (this.get('controller').get('postButtonDisabled')) {
        return;
      }

      // Ctrl + Enter.
      if (event.ctrlKey && event.which === 13) {
        this.get('controller').send('createNewPost');
      }
    }
  });

  //-------------------------
  // Controllers
  //-------------------------
  App.ApplicationController = Ember.Controller.extend({
    needs: ['posts', 'users']
  });

  App.IndexController = Ember.Controller.extend({
    needs: ['posts', 'users'],
    myUserBinding: 'controllers.users.myUser',
    newPostTitle: '',

    postBoxDisabled: function() {
      return !(this.get('myUser.isLoaded') && (this.get('myUser.length') > 0));
    }.property('myUser.isLoaded', 'myUser.length'),

    postButtonDisabled: function() {
      return !(!this.get('postBoxDisabled') && this.get('newPostTitle.length') > 0);
    }.property('postBoxDisabled', 'newPostTitle'),

    createNewPost: function() {
      var self = this;
      var myUser = this.get('controllers.users.myUser.firstObject.user');

      var newPost = App.Post.createRecord({
        user: myUser,
        title: this.get('newPostTitle')
      });
      newPost.on('didCreate', function() {
        self.set('newPostTitle', '');
      });
      newPost.get('transaction').commit();
    }
  });

  App.ProfileController = Ember.Controller.extend({
    needs: ['posts', 'users'],
    myUserBinding: 'controllers.users.myUser',
    filterByUserBinding: 'controllers.posts.filterByUser',
    newName: '',
    newEmail: '',

    // Helper to disable the button when the fields are not filled.
    createButtonDisabled: function() {
      return this.get('newName.length') === 0;
    }.property('newName'),

    createNew: function() {
      var self = this;

      var newUser = App.User.createRecord({
        name: this.get('newName'),
        email: this.get('newEmail'),
        isMyUser: true
      });

      newUser.on('didCreate', function() {
        self.set('myUser', App.User.filter(function(user) {
          return user.get('isMyUser') === true;
        }));
        self.set('filterByUser', self.get('myUser'));
      });

      newUser.get('transaction').commit();
    }
  });

  App.UsersController = Ember.ArrayController.extend({
    // This is set to a FilteredRecordArray by the router. Just use the
    // first object in the array.
    myUser: []
  });

  App.PostsController = Ember.ArrayController.extend({
    filterByUser: [],
    sortProperties: ['created'],
    sortAscending: false,

    filteredPosts: function() {
      if (this.get('filterByUser.length') > 0) {
        var filterUserId = this.get('filterByUser.firstObject.id');
        if (filterUserId) {
          return this.get('arrangedContent').filterProperty('userId', filterUserId);
        }
      }
      return this.get('arrangedContent');
    }.property('content.[]', 'filterByUser.[]'),

    deletePost: function(id) {
      if (confirm('Are you sure you want to delete this post?')) {
        var post = App.Post.find(id);
        post.deleteRecord();
        post.get('transaction').commit();
      }
    },

    repostPost: function(id) {
      if (confirm('Are you sure you want to repost this?')) {
        var post = App.Post.find(id);
        var myUser = this.get('controllers.users.myUser.firstObject.user');
        var newPost = App.Post.createRecord({
            user: myUser,
            title: post.get('title'),
            repostUserName: post.get('userName'),
            repostUserAvatar: post.get('userAvatar'),
            repostPostId: post.get('id')
          });
          newPost.get('transaction').commit();
      }
    },

    loadMore: function() {
      App.Post.find({before : this.get('filteredPosts.lastObject.id')});
    }
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
      controller.set('controllers.users.myUser', App.User.find({'is_my_user': true}));
      var refreshUI = function() {
        App.Post.find();
        setTimeout(refreshUI, Config.ui.pollInterval);
      };
      setTimeout(refreshUI, 5000);
    }
  });

  App.IndexRoute = Ember.Route.extend({
    setupController: function(controller) {
      var postsController = controller.get('controllers.posts');
      postsController.set('filterByUser', []);
    }
  });

  App.ProfileRoute = Ember.Route.extend({
    setupController: function(controller) {
      var postsController = controller.get('controllers.posts');
      var usersController = controller.get('controllers.users');
      postsController.set('filterByUser', usersController.get('myUser'));
    }
  });

  // TODO: Put this somewhere else.
  $('.time').moment({ frequency: 5000 });

});
