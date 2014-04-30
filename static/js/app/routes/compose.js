define(function(require) {

  var $ = require('jquery');
  var ComposeView = require('app/views/compose');

  return function(id) {
    var view = new ComposeView({
      id: id
    });

    vole.view.setContentView(view);
  };

});
