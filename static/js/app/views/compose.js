define(function(require) {

  var Backbone = require('backbone');
  var Handlebars = require('handlebars');

  var Post = require('app/models/post');
  var Drafts = require('app/collections/draft');

  var EditorView = require('app/views/editor');
  var DraftsView = require('app/views/drafts');

  return Backbone.View.extend({

    className: 'compose',

    template: Handlebars.compile(require('text!tmpl/compose.hbs')),

    initialize: function(options) {
      this.options = options;
    },

    render: function() {
      this.$el.html(this.template());

      var post = new Post();

      if (this.options.id) {
        post.set('draft', true);
        post.set('id', this.options.id);
        post.fetch();
      }

      this.$('#editor').html(new EditorView({
        model: post
      }).render().el);

      this.$('#drafts').html(new DraftsView({
        collection: new Drafts()
      }).render().el);

      return this;
    }

  });

});
