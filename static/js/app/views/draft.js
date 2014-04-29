define(function(require) {

  var Backbone = require('backbone');
  var Handlebars = require('handlebars');
  var marked = require('lib/marked');

  return Backbone.View.extend({

    template: Handlebars.compile(require('text!tmpl/draft.hbs')),

    summary: function() {
      var html = marked(this.model.get('title'));
      return $('<div>').html(html).text().trim().substring(0, 30);
    },

    render: function() {
      this.model.set('summary', this.summary());
      this.$el.html(this.template(this.model.attributes));
      return this;
    }

  });

});
