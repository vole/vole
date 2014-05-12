
requirejs.config({
  paths : {
    'moment': 'lib/moment',
    'text': 'plugins/text',
    'tmpl': 'app/templates',
    'underscore': 'lib/underscore',
    'backbone': 'lib/backbone',
    'handlebars': 'lib/handlebars',
    'foundation': 'lib/foundation'
  },
  shim: {
    'handlebars': {
      exports: 'Handlebars'
    },
    'foundation': {
      exports: 'Foundation'
    }
  }
});

define(function(require) {

  // So we can use es5 features.
  var es5shim = require('lib/es5-shim');
  var es5sham = require('lib/es5-sham');

  // Make sure handlebars helpers have loaded before we try
  // to render any templates.
  var helpers = require('tmpl/helpers');

  var _ = require('underscore');
  var Backbone = require('backbone');;

  var ConfigModel = require('app/models/config');
  var UserModel = require('app/models/user');
  var AppModel = require('app/models/app');

  var AppView = require('app/views/app');

  var Router = require('app/router');

  vole = {};

  vole.logger = require('lib/logger')('vole');

  // Create a global event bus.
  vole.events = _.extend({}, Backbone.Events);

  // The main view that contains the entire application.
  vole.view = new AppView({
    model: new AppModel(),
    el: '#app'
  });

  // Create the main app router.
  vole.router = new Router();

  // Initialize the app config.
  vole.config = new ConfigModel();

  // Represents the current user.
  vole.user = new UserModel();

  vole.logger.info('starting');
  vole.logger.info('loading config');

  // Once the config has loaded, we can load the user's information.
  //
  // TODO: If the user hasn't been created, we need to kick off the
  // user creation flow.
  vole.config.on('sync', function() {
    vole.logger.info('config loaded');

    vole.view.render();

    // Fetch the current user.
    vole.user.fetch();
  });

  vole.logger.info('loading user');

  // At this point, the configuration and the user have been fully
  // initialized, and it's safe to start the application.
  vole.user.on('sync', function() {
    vole.logger.info('user loaded');
    vole.logger.info('starting router');

    // Start the main app router.
    Backbone.history.start({ pushState: true });
  });

  vole.user.on('error', function(user, response, options) {
    vole.logger.info('user not found');
    vole.logger.info('starting router');

    // Start the main app router.
    Backbone.history.start({ pushState: true });
    Backbone.history.navigate('install', { trigger: true });
  });

  vole.router.on('route', function(route, params) {
    vole.logger.info('route:', route, 'params:', params.join(', '));
  });

  // Fetch the app config.
  vole.config.fetch();
});
