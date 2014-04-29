define(function(require) {

  var Backbone = require('backbone');
  var Handlebars = require('handlebars');
  var DraftView = require('app/views/draft')

  return Backbone.View.extend({

    className: 'drafts',

    template: Handlebars.compile(require('text!tmpl/drafts.hbs')),

    initialize: function() {
      this.collection.on('sync', this.render.bind(this));
      this.collection.fetch();
    },

    render: function() {
      this.$el.html(this.template());

      this.collection.each(function(draft) {
        this.$('ul').append(
          new DraftView({
            model: draft,
            tagName: 'li'
          }).render().el
        );
      }.bind(this));

      return this;
    }

  });

});
