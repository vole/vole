define(function(require) {

  var _ = require('underscore');

  var formatRegExp = /%[sdj%]/g;

  /**
   * Node-style string formatting.
   *
   * From: https://github.com/joyent/node/blob/master/lib/util.js#L22
   *
   * @param  {String} f
   * @return {String}
   */
  return function(f) {
    var i;

    if (!_.isString(f)) {
      var objects = [];
      for (i = 0; i < arguments.length; i++) {
        objects.push(inspect(arguments[i]));
      }
      return objects.join(' ');
    }

    i = 1;
    var args = arguments;
    var len = args.length;
    var str = String(f).replace(formatRegExp, function(x) {
      if (x === '%%') return '%';
      if (i >= len) return x;
      switch (x) {
        case '%s': return String(args[i++]);
        case '%d': return Number(args[i++]);
        case '%j':
          try {
            return JSON.stringify(args[i++]);
          }
          catch (_) {
            return '[Circular]';
          }
          break;
        default:
          return x;
      }
    });

    return str;
  };

});
