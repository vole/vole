define(function(require) {

  var Backbone = require('backbone');
  var Handlebars = require('handlebars');

  var BaseView = require('app/views/base');
  var DraftView = require('app/views/draft')

  return BaseView.extend({

    className: 'drafts',

    template: Handlebars.compile(require('text!tmpl/drafts.hbs')),

    initialize: function() {
      this.collection.on('sync', this.render.bind(this));
      this.collection.fetch();
    },

    render: function() {
      this.$el.html(this.template());

      this.collection.each(function(draft) {
        this.subView('ul', new DraftView({
          model: draft,
          tagName: 'li'
        }));
      }.bind(this));

      return this;
    }

  });

});
