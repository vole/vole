define(function(require) {

  var underscore = require('underscore');
  var Backbone = require('backbone');
  var Handlebars = require('handlebars');

  return Backbone.View.extend({

    subView: function(selector, view) {
      if (!this._subViews) {
        this._subViews = {};
      }

      this._subViews[selector] = view;

      var $el = view.render().$el;

      $el.appendTo(this.$(selector));
    }

  });

  // var ModelView = Backbone.View.extend({

  //   // Override this with a real template.
  //   template: '',

  //   initialize: function() {
  //     _.extend(this, Backbone.Events);

  //     this._subViews = {};

  //     if (this.afterInitialize) {
  //       this.afterInitialize();
  //     }
  //   },

  //   subView: function(selector, view) {
  //     this._subViews[selector] = view;

  //     // Bubble events up to this view.
  //     view.on('all', this.trigger.bind(this));
  //   },

  //   render: function() {
  //     // Compile the template.
  //     var compiledTemplate = Handlebars.compile(this.template);

  //     // Render the HTML.
  //     var html = compiledTemplate(this.model.attributes);

  //     // Set the view's HTML.
  //     this.$el.html(html);

  //     // Render sub views.
  //     this.renderSubViews();

  //     return this;
  //   },

  //   renderSubView: function(selector) {
  //     var view = this._subViews[selector].render();
  //     this.$(selector).html(view.el);
  //   },

  //   renderSubViews: function() {
  //     Object.keys(this._subViews).forEach(this.renderSubView.bind(this));
  //   }

  // });

  // var CollectionView = ModelView.extend({

  //   render: function() {
  //     this.collection.each(function(model) {
  //       var view = new this.View({
  //         model: model
  //       });

  //       view.on('all', this.trigger.bind(this));

  //       this.$el.append(view.el);
  //     }.bind(this));
  //   }

  // });

  // return {
  //   ModelView: ModelView,
  //   CollectionView: CollectionView
  // };

  // var PostView = ModelView.extend({

  // });

  // var PostsView = CollectionView.extend({
  //   View: PostView
  // });

  // var TimeLineView = ModelView.extend({

  //   afterInitialize: function() {
  //     var posts = new PostsView({
  //       collection: new PostsCollection()
  //     });

  //     this.subView('.posts', posts);
  //   }

  // });

  // new TimeLineView().render();

});
