
(function ($) {

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

})(jQuery);
