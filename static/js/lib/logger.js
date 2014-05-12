define(function(require) {

  var levels = [
    'off',
    'fatal',
    'error',
    'warn',
    'info',
    'debug',
    'trace'
  ];

  var toArray = function(obj) {
    return Array.prototype.slice.call(obj, 0);
  };

  function Logger(name) {
    this.name = name;
  }

  Logger.prototype.log = function(level) {
    level = levels.indexOf(level);
    var threshold = levels.indexOf(vole.config.get('ui_logging'));

    if (level > threshold) {
      return;
    }

    var args = Array.prototype.slice.call(arguments, 1);
    var fn = console.log;

    if (level === levels.indexOf('info')) {
      fn = console.info;
    }
    else if (level === levels.indexOf('warn')) {
      fn = console.warn;
    }
    else if (level > levels.indexOf('warn')) {
      fn = console.error;
    }

    fn.apply(console, ['[' + this.name + ']'].concat(args));
  };

  levels.forEach(function(level) {
    Logger.prototype[level] = function() {
      this.log.apply(this, [level].concat(toArray(arguments)));
    };
  });

  return function(name) {
    return new Logger(name);
  };

});
