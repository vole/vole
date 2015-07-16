
requirejs.config({
  paths : {
    'moment': 'lib/moment',
    'text': 'plugins/text',
    'tmpl': 'app/templates',
    'underscore': 'lib/underscore',
    'backbone': 'lib/backbone',
    'handlebars': 'lib/handlebars',
    'foundation': 'lib/foundation',
    'highlight': 'lib/highlight'
  },
  shim: {
    'handlebars': {
      exports: 'Handlebars'
    },
    'foundation': {
      exports: 'Foundation'
    },
    'highlight': {
      exports: 'hljs'
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

  var Cache = require('lib/cache').Memory;

  // Global namespace.
  vole = {};

  // Catch-all logger, modules should use their own scoped logger.
  vole.logger = require('lib/logger')('vole');

  // Create a global event bus. Eventually this can be used to
  // communicate with plugins.
  vole.events = _.extend({}, Backbone.Events);

  // The main view that contains the entire application.
  vole.view = new AppView({
    model: new AppModel(),
    el: '#app'
  });

  vole.cache = new Cache();

  // Create the main app router.
  vole.router = new Router();

  // Initialize the app config.
  vole.config = new ConfigModel();

  // Represents the current user.
  vole.user = new UserModel();

  // Log router events.
  vole.router.on('route', function(route, params) {
    vole.logger.info('route:', route, 'params:', params.join(', '));
  });

  // This is important.
  console.log(
    '\n' +
    ' _    __      __\n' +
    '| |  / /___  / /__\n' +
    '| | / / __ \\/ / _ \\\n' +
    '| |/ / /_/ / /  __/\n' +
    '|___/\\____/_/\\___/\n\n'
  );

  vole.logger.info('loading config');

  // Once the config has loaded, we can load the user's information.
  vole.config.on('sync', function() {
    vole.logger.info('config', vole.config.attributes);

    // Once we have the config we can safely render the app view.
    vole.view.render();

    vole.logger.info('loading user');

    // Fetch the current user.
    vole.user.fetch();
  });

  // At this point, the configuration and the user have been fully
  // initialized, and it's safe to start the application.
  vole.user.on('sync', function() {
    vole.logger.info('user', vole.user.attributes);
    vole.logger.info('starting router');

    // Start the main app router.
    Backbone.history.start({ pushState: true });
  });

  // If there is an error loading the user, it means we need to kick
  // off the installation process.
  vole.user.on('error', function(user, response, options) {
    vole.logger.warn('user not found');
    vole.logger.info('starting router');

    // Start the main app router.
    Backbone.history.start({ pushState: true });
    Backbone.history.navigate('install', { trigger: true });
  });

  // Fetch the app config.
  vole.config.fetch();
});
