
module.exports = function(grunt) {

  var util = require('util')
  var shelljs = require('shelljs')
  var mkdirp = require('mkdirp')

  grunt.registerTask('macgap', function() {
    grunt.log.subhead('macgap')

    var config = grunt.config(this.name)
    var source = config.src
    var destination = config.dest

    mkdirp.sync(destination)

    var command = util.format('macgap build -o %s %s', destination, source)

    grunt.log.write(command)

    var result = shelljs.exec(command)

    if (result.code === 0) {
      grunt.log.ok()
    }
    else {
      grunt.log.error(result.output)
      return false
    }
  });

};
