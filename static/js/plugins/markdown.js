
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

    var html = '';

    try {
      html = marked(content || '');
    }
    catch (e) {
      html = '[Parse Error]';
    }

    return new Handlebars.SafeString(html);
  });

});
