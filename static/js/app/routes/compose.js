define(function(require) {

  var $ = require('jquery');
  var ComposeView = require('app/views/compose');

  return function(id) {
    $('#content').html(new ComposeView({
      id: id
    }).render().el);
  };

});
