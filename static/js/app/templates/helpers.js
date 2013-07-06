
define([
  'ember',
  'handlebars',
  'lib/marked',
  'moment',
  'jquery'
],
function (Ember, Handlebars, marked, moment, $) {

  Ember.Handlebars.registerBoundHelper('nanoDate', function(value, options) {
    var escaped = Handlebars.Utils.escapeExpression(value);
    var ms = Math.round(escaped / Math.pow(10, 6));
    return new Handlebars.SafeString(moment(ms).fromNow());
  });

  Ember.Handlebars.registerBoundHelper('markdown', function(content) {
    marked.setOptions({
      gfm: true,
      sanitize: true,
      breaks: true
    });

  var html = $('<div>')
    .html(marked(content))
    .find('a')
      .attr('target', '_blank')
    .end()
    .html();

    return new Handlebars.SafeString(html);
  });

});
