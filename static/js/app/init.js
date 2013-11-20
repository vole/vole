
define([
  'app/timeline',
  'app/header',
  'app/postbox',
  'app/sidebar',
  'app/create_profile',
  'app/templates/helpers/all'
],
function (Timeline, Header, Postbox, Sidebar, CreateProfile) {
	Timeline.attachTo('.timeline');
	Header.attachTo('#main-nav');
	Postbox.attachTo('#create-post');
	Sidebar.attachTo('#sidebar');
	CreateProfile.attachTo('#create-profile');
});
