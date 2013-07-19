
define([
  'app/config',
  'ember',
  'ember-data',
  'lib/marked',
  'lib/socket.io-client',
  'plugins/text!app/templates/application.hbs',
  'plugins/text!app/templates/index.hbs',
  'plugins/text!app/templates/posts.hbs',
  'plugins/text!app/templates/profile.hbs',
  'app/templates/helpers',
  'plugins/moment',
  'plugins/resize'
],
function (Config, Ember, DS, marked, io, applicationTemplate, indexTemplate, postsTemplate, profileTemplate) {
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
    isMyPost: DS.attr('boolean')
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
    }.property('content.[]', 'filterByUser.[]', 'arrangedContent.[]'),

    deletePost: function(id) {
      if (confirm('Are you sure you want to delete this post?')) {
        var post = App.Post.find(id);
        post.deleteRecord();
        post.get('transaction').commit();
      }
    },

    loadMore: function() {
      var filter = {
        'before': this.get('filteredPosts.lastObject.id')
      };
      var filterUserId = this.get('filterByUser.firstObject.id');
      if (filterUserId) {
        filter.user = filterUserId;
      }
      App.Post.find(filter);
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
      postsController.set('content', App.Post.find());
    }
  });

  App.ProfileRoute = Ember.Route.extend({
    setupController: function(controller) {
      var postsController = controller.get('controllers.posts');
      var usersController = controller.get('controllers.users');
      postsController.set('filterByUser', usersController.get('myUser'));
      App.Post.find({user: 'my_user'});
      postsController.set('content', App.Post.filter(function(item) {
        return item.get('isMyPost');
      }));
    }
  });

  // TODO: Put this somewhere else.
  $('.time').moment({ frequency: 5000 });

  //-------------------------
  // Websockets
  //-------------------------
  /*
  var socket = io.connect('ws://127.0.0.1:6789', {
    resource: 'ws'
  });
  socket.on('connect', function() {
    console.log('boom');
  });
  */
  var conn = new WebSocket("ws://localhost:6789/ws");
  conn.onopen = function(evt) {
    console.log('Connection opened.');
  };
  conn.onclose = function(evt) {
    console.log('Connection closed.');
  };
  conn.onmessage = function(evt) {
    console.log(evt.data);
  };
});
