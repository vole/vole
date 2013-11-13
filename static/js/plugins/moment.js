
define(['moment', 'jquery'], function (moment, $) {

  moment.lang('en', {
    relativeTime : {
      future: "in %s",
      past:   "%s",
      s:  "%d s",
      m:  "1 m",
      mm: "%d m",
      h:  "1 h",
      hh: "%d h",
      d:  "1 d",
      dd: "%d d",
      M:  "1 mon",
      MM: "%d mon",
      y:  "1 y",
      yy: "%d y"
    }
  });

  $.fn.moment = function (options) {
    options = options || {};

    var frequency = options.frequency || 1000;
    var selector = this.selector;

    var poll = function () {
      $(selector).each(function () {
        var $this = $(this);
        var ts = $this.data('time');
        var ms = Math.round(parseInt(ts, 10) / Math.pow(10, 6));
        $this.text(moment(ms).fromNow());
      });

      setTimeout(poll, frequency);
    };

    poll();
  };

});
