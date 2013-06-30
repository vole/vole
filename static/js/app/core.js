
define([
  'app/config',
  'ember',
  'ember-data',
  'lib/showdown',
  'plugins/text!app/templates/application.hbs',
  'plugins/text!app/templates/index.hbs',
  'plugins/text!app/templates/posts.hbs',
  'plugins/text!app/templates/profile.hbs',
  'plugins/moment',
  'plugins/resize'
],
function (Config, Ember, DS, showdown, applicationTemplate, indexTemplate, postsTemplate, profileTemplate) {

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

  //-------------------------
  // Handlebars
  //-------------------------
  Ember.Handlebars.registerBoundHelper('nanoDate', function(value, options) {
    var escaped = Handlebars.Utils.escapeExpression(value);
    var ms = Math.round(escaped / Math.pow(10, 6));
    return new Handlebars.SafeString(moment(ms).fromNow());
  });

  var linkify = function (text) {
    var rUrl = /\(?\b(?:(http|https|ftp):\/\/)+((?:www.)?[a-zA-Z0-9\-\.]+[\.][a-zA-Z]{2,4}|localhost(?=\/)|\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})(?::(\d*))?(?=[\s\/,\.\)])([\/]{1}[^\s\?]*[\/]{1})*(?:\/?([^\s\n\?\[\]\{\}\#]*(?:(?=\.)){1}|[^\s\n\?\[\]\{\}\.\#]*)?([\.]{1}[^\s\?\#]*)?)?(?:\?{1}([^\s\n\#\[\]\(\)]*))?([\#][^\s\n]*)?\)?/ig;
    var matches = text.match(rUrl);

    if (matches) {
      for (var x = 0, len = matches.length; x < len; x++) {
        var match = matches[x];

        var outer = $('<div />');
        var link = $('<a />', {
          href : match,
          target : '_blank'
        });

        if (/\.(jpg|jpeg|gif|png|bmp|ico)$/.test(match)) {
          var image = $('<img />', { src : match });
          image.addClass('img-rounded');
          link.html(image);
        }
        else {
          link.text(match);
        }

        text = text.replace(match, outer.append(link).html());
      }
    }

    return text;
  };

  Ember.Handlebars.registerBoundHelper('markdown', function(content) {
    // Convert the markdown to HTML.
    var converter = new Showdown.converter();
    var converted = converter.makeHtml(content);

    // Create a fragment with the HTML so we can traverse its text nodes.
    var html = $('<div>').html(converted);

    function collectTextNodes (element, texts) {
      for (var child = element.firstChild; child !== null; child = child.nextSibling) {
        if (child.nodeType === 3) {
          texts.push(child);
        }
        else if (child.nodeType === 1) {
          collectTextNodes(child, texts);
        }
      }
    }

    // Find all the text nodes.
    var textNodes = [];
    collectTextNodes(html.get(0), textNodes);

    // Replace each text node with a linkified version.
    for (var x = 0; x < textNodes.length; x++) {
      $(textNodes[x]).replaceWith(linkify(textNodes[x].data));
    }

    // Return the complete HTML.
    return new Handlebars.SafeString(html.html());
  });

  $('.time').moment({ frequency: 5000 });

});
