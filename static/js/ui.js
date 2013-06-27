/**
 * Pretty stuff
 */
(function($){
  // Expand images to full-size clicked
  $(document).on('click', '.post .span6 img', function(e) {
    e.preventDefault();
    $(this).toggleClass('full');
  });
})(jQuery);
