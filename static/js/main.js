
requirejs.config({
  paths : {
    'moment' : 'lib/moment',
    'flight' : 'lib/flight',
    'text' : 'plugins/text',
    'tmpl' : 'app/templates',
    'bootstrap': 'lib/bootstrap'
  }
});

define([
  // Global dependencies.
  'lib/es5-shim',
  'lib/es5-sham',
  'plugins/markdown',
  'plugins/nanodate',
  'plugins/moment',
  'bootstrap'
], function () {
  // Initialize app.
  require(['app/init']);
});
