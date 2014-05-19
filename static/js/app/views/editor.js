define(function(require) {

  require('lib/jquery.autosize');

  var _ = require('underscore');
  var BaseView = require('app/views/base');
  var Handlebars = require('handlebars');
  var marked = require('lib/marked');
  var highlight = require('highlight');
  var Post = require('app/models/post');

  marked.setOptions({
    sanitize: true,
    gfm: true,
    highlight: function (code) {
      return highlight.highlightAuto(code).value;
    }
  });

  return BaseView.extend({

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

      this.showLoading();

      this.model.save({
        title: body
      }, {
        // On successful save, redirect to the draft view. This intuitively
        // updates the route and sidebar.
        success: function(draft) {
          this.hideLoading();
          Backbone.history.navigate('/compose/' + draft.get('id'), true);
        }.bind(this),

        // TODO(aaron): Handle draft save error.
        error: function() {
          this.hideLoading();
        }.bind(this)
      });
    },

    // TODO(aaron): Handle errors.
    post: function(e) {
      e.preventDefault();

      var body = this.body();
      if (!body) {
        return;
      }

      var post = new Post({
        title: body
      });

      post.save();

      this.model.destroy();

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
