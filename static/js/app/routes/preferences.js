define(function(require) {

  var $ = require('jquery');
  var PreferencesView = require('app/views/preferences');

  return function(action) {
    $('#content').html(new PreferencesView().render().el);
  };

});
