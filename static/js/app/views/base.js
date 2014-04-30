define(function(require) {

  var Backbone = require('backbone');

  return Backbone.View.extend({

    // Default parent reference so views don't have to
    // constantly check if they are a child or not.
    //parent: new Backbone.View(),

    // Add a sub-view to this view.
    subView: function(selector, view) {
      if (!this._subViews) {
        this._subViews = [];
      }

      // Create references between parent and child.
      this._subViews.push(view);
      view.parent = this;

      // Render the view and attach it to the specified selector.
      var $el = view.render().$el;
      $el.appendTo(this.$(selector));
    }

  });

});
