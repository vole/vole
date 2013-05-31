module.exports = function(grunt) {
  grunt.initConfig({
    clean: ['dist/', 'bin/vole'],

    copy: {
      dist: {
        files: [
          {
            src: ['static/**'],
            dest: 'dist/'
          },
          {
            src: ['bin/vole'],
            dest: 'dist/vole'
          },
          {
            src: ['bin/vole.exe'],
            dest: 'dist/vole.exe'
          },
          {
            src: ['config.sample.json'],
            dest: 'dist/config.sample.json'
          }
        ]
      }
    },

    chmod: {
      options: {
        mode: '755'
      },
      vole: {
        src: ['dist/vole']
      }
    },

    exec: {
      install_vole: {
        command: 'go install vole'
      }
    },

    jshint: {
      options: {
        scripturl: true
      },
      prod: {
        files: {
          src: ['Gruntfile.js', 'static/js/*.js']
        }
      },
      dev: {
        options: {
          debug: true
        },
        files: {
          src: ['Gruntfile.js', 'static/js/*.js']
        }
      }
    }
  });

  grunt.loadNpmTasks('grunt-contrib-jshint');
  grunt.loadNpmTasks('grunt-contrib-clean');
  grunt.loadNpmTasks('grunt-contrib-copy');
  grunt.loadNpmTasks('grunt-chmod');
  grunt.loadNpmTasks('grunt-exec');

  grunt.registerTask('default', ['jshint', 'clean', 'exec:install_vole', 'copy', 'chmod']);
};
