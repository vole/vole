define(function(require) {

  var $ = require('jquery');
  var PreferencesView = require('app/views/preferences');

  return function(action) {
    var view = new PreferencesView();

    vole.view.setContentView(view);
  };

});
