(function ($) {
  var cl = console.log.bind(console);
  var App = Ember.Application.create({
    LOG_TRANSITIONS: true,
    rootElement: '#ember-container'
  });
  window.App = App;

  //-------------------------
  // Models
  //-------------------------

  //-------------------------
  // Views
  //-------------------------

  //-------------------------
  // Controllers
  //-------------------------

  //-------------------------
  // Router
  //-------------------------

  App.Router.map(function() {
    this.resource("index", {path: "/"});
  });

})(jQuery);
