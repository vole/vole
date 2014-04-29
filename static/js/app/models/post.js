define(function(require) {

  var Backbone = require('backbone');

  return Backbone.Model.extend({

    url: function() {
      if (this.get('draft')) {
        return '/api/drafts/' + this.get('id');
      }

      return '/api/posts' + (this.isNew() ? '' : '/' + this.get('id'));
    }

  });

});
