
define([
  'app/config',
  'app/timeline',
  'app/header',
  'app/postbox',
  'app/sidebar',
  'app/create_profile'
],
function (Config, Timeline, Header, Postbox, Sidebar, CreateProfile) {
	Timeline.attachTo('.timeline');
	Header.attachTo('#main-nav');
	Postbox.attachTo('#above-posts');
	Sidebar.attachTo('#sidebar');
	CreateProfile.attachTo('#create-profile');
});
