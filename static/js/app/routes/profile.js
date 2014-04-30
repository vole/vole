define(function(require) {

  var $ = require('jquery');
  var ProfileView = require('app/views/profile');

  return function(key) {
    var view = new ProfileView({
      model: vole.user
    });

    vole.view.setContentView(view);
  };

});
