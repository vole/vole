module.exports = function(grunt) {
  grunt.initConfig({
    clean: ['dist/', 'bin/vole', 'bin/vole.exe'],

    copy: {
      dist: {
        files: [
          {
            src: ['bin/vole'],
            dest: 'dist/vole'
          },
          {
            src: ['bin/vole.exe'],
            dest: 'dist/vole.exe'
          },
          {
            src: [
              'static/**',
              'README.md',
              'CONTRIBUTING.md',
              'CHANGELOG.md',
              'LICENSE',
              'config.sample.json'
            ],
            dest: 'dist/'
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
    },

    macgap: {
      src: './build/osx/vole',
      dest: './dist/'
    }
  });

  grunt.loadNpmTasks('grunt-contrib-jshint');
  grunt.loadNpmTasks('grunt-contrib-clean');
  grunt.loadNpmTasks('grunt-contrib-copy');
  grunt.loadNpmTasks('grunt-chmod');
  grunt.loadNpmTasks('grunt-exec');

  grunt.loadTasks('./build/tasks');

  grunt.registerTask('default', ['jshint', 'clean', 'exec:install_vole', 'copy', 'chmod']);

  grunt.registerTask('build:osx', ['macgap']);

  grunt.registerTask('build', ['clean', 'build:osx']);
};
