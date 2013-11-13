
define([
  'lib/handlebars',
  'lib/marked'
],
function (Handlebars, marked) {

  Handlebars.registerHelper('markdown', function (content) {
    marked.setOptions({
      gfm: true,
      sanitize: true,
      breaks: true
    });

    return new Handlebars.SafeString(marked(content));
  });

});
