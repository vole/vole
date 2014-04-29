define(function(require) {

  var Handlebars = require('handlebars');
  var marked = require('lib/marked');
  var moment = require('moment');

  Handlebars.registerHelper('nanoDate', function(value, options) {
    var escaped = Handlebars.Utils.escapeExpression(value);
    var ms = Math.round(escaped / Math.pow(10, 6));
    return new Handlebars.SafeString(moment(ms).fromNow());
  });

  Handlebars.registerHelper('markdown', function(content) {
    marked.setOptions({
      gfm: true,
      sanitize: true,
      breaks: true
    });

    var html = '';

    try {
      html = marked(content);
    }
    catch (e) {
      console.error(e.stack);
      html = '[Parse Error]';
    }

    return new Handlebars.SafeString(html);
  });

  Handlebars.registerHelper('equal', function(a, b, options) {
    if (a == b) {
      return options.fn(this);
    }
    else {
      return options.inverse(this);
    }
  });

});
