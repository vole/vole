
define([
  'lib/handlebars',
  'moment'
],
function (Handlebars, moment) {

  Handlebars.registerHelper('nanoDate', function (value, options) {
    var escaped = Handlebars.Utils.escapeExpression(value);
    var ms = Math.round(escaped / Math.pow(10, 6));
    return new Handlebars.SafeString(moment(ms).fromNow());
  });

});
