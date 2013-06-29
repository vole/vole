
requirejs.config({
  paths : {
    'ember' : 'lib/ember/ember',
    'ember-data' : 'lib/ember/ember-data',
    'handlebars' : 'lib/handlebars',
    'moment' : 'lib/moment'
  },
  shim : {
    'ember' : {
      deps : ['jquery', 'handlebars'],
      exports : 'Ember'
    },
    'ember-data' : {
      deps : ['ember'],
      exports : 'DS'
    }
  }
});

define(['app/core'], function () {});
