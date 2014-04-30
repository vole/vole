define(function(require) {

  var Backbone = require('backbone');
  var Handlebars = require('handlebars');

  var Post = require('app/models/post');
  var Drafts = require('app/collections/draft');

  var BaseView = require('app/views/base');
  var EditorView = require('app/views/editor');
  var DraftsView = require('app/views/drafts');

  return BaseView.extend({

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

      this.subView('#editor', new EditorView({
        model: post
      }));

      this.subView('#drafts', new DraftsView({
        collection: new Drafts()
      }));

      return this;
    }

  });

});
