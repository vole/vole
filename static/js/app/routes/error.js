define(function(require) {

  var $ = require('jquery');
  var ErrorView = require('app/views/error');

  return function(id) {
    var view = new ErrorView();
    vole.view.setContentView(view);
  };

});
