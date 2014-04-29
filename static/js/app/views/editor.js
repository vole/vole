define(function(require) {

  require('lib/jquery.autosize');

  var _ = require('underscore');
  var Backbone = require('backbone');
  var Handlebars = require('handlebars');
  var marked = require('lib/marked');

  return Backbone.View.extend({

    className: 'editor',

    template: Handlebars.compile(require('text!tmpl/editor.hbs')),

    events: {
      'keyup .js-editor': 'preview',
      'click .js-post': 'post',
      'click .js-save': 'save'
    },

    initialize: function() {
      this.model.on('sync', this.render.bind(this));
    },

    body: function() {
      return this.$('.js-editor').val();
    },

    preview: _.debounce(function() {
      this.$('.js-preview').html(marked(this.body()));
    }, 100),

    save: function(e) {
      e.preventDefault();

      var body = this.body();

      if (!body) {
        return;
      }

      this.model.set('draft', true);
      this.model.set('title', body);
      this.model.save();
    },

    post: function(e) {
      e.preventDefault();

      var body = this.body();

      if (!body) {
        return;
      }

      this.model.set('draft', false);
      this.model.set('title', body);
      this.model.save();

      Backbone.history.navigate('/timeline', true);
    },

    render: function() {
      this.$el.html(this.template(this.model.attributes));
      this.preview();
      this.$('textarea').autosize();
      return this;
    }

  });

});
