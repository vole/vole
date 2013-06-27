/**
 * Pretty stuff
 */
(function($){
  // Expand images to full-size clicked
  $('.post img').click(function(e){
    e.preventDefault();
    $(this).toggleClass('full');
  });
})(jQuery);
