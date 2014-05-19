define(function(require) {

  var md5 = require('lib/md5');
  var format = require('lib/format');

  return function(email) {
    var hash = md5(email.trim().toLowerCase());
    return format('http://www.gravatar.com/avatar/%s', hash);
  };

});
