define(['lib/handlebars'], function (Handlebars) {

	Handlebars.registerHelper('formatFileSize', function (bytes, si) {
		var threshold = si ? 1000 : 1024;

		if (bytes < threshold) {
			return bytes + ' B';
		}

		var units = si ?
			['kB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'] :
			['KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];

		var u = -1;

		do {
			bytes /= threshold;
			++u;
		} while (bytes >= threshold);

		return bytes.toFixed(1) + ' ' + units[u];
	});

});
