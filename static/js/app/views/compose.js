define(function(require) {

  var Backbone = require('backbone');
  var Handlebars = require('handlebars');

  var Draft = require('app/models/draft');
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

      var draft = new Draft();

      if (this.options.id) {
        draft.set('id', this.options.id);
        draft.fetch();
      }

      this.subView('#editor', new EditorView({
        model: draft
      }));

      this.subView('#drafts', new DraftsView({
        collection: new Drafts()
      }));

      return this;
    }

  });

});
