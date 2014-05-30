define(function(require) {

  var TimelineView = require('app/views/timeline');

  return function(user) {
    var view = new TimelineView({
      user: user
    });

    vole.view.setContentView(view);
  };

});
