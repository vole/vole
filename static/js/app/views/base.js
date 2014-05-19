define(function(require) {

  var Backbone = require('backbone');
  var Spin = require('lib/spin');

  return Backbone.View.extend({

    /**
     * Render a subview.
     *
     * @param  {String}        selector Parent element to append rendered
     *                                  element to.
     * @param  {Backbone.View} view     The Backbone view.
     * @return {Backbone.View}          this
     */
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

      return this;
    },

    subViews: function() {
      return this._subViews;
    },

    /**
     * Show a loading overlay over the view.
     *
     * TODO(aaron): Make this more effecient!
     */
    showLoading: function() {
      var position = this.$el.offset();
      var width = this.$el.outerWidth();
      var height = this.$el.outerHeight();

      this.overlay = $('<div>', { 'class': 'loading' }).css({
        top: position.top + 'px',
        left: position.left + 'px',
        width: width + 'px',
        height: height + 'px'
      }).appendTo('body');

      var spinner = new Spin({
        width: 2
      }).spin(this.overlay.get(0));

      return this;
    },

    /**
     * Hide the loading overlay.
     */
    hideLoading: function() {
      if (this.overlay) {
        this.overlay.remove();
      }
      return this;
    }

  });

});
