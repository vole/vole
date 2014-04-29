define(function(require) {

  var $ = require('jquery');
  var TimelineView = require('app/views/timeline');

  return function(user) {
    var timeline = new TimelineView({ user: user });
    $('#content').html(timeline.render().el);
  };

});
