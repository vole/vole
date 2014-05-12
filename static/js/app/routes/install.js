define(function(require) {

  var InstallView = require('app/views/install');
  var InstallModel = require('app/models/install');

  return function(id) {
    var view = new InstallView({
      model: new InstallModel()
    });

    vole.view.setContentView(view);
  };

});
