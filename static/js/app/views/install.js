define(function(require) {

  var Handlebars = require('handlebars');
  var Backbone = require('backbone');
  var Spin = require('lib/spin');
  var logger = require('lib/logger')('install');
  var User = require('app/models/user');

  return Backbone.View.extend({

    className: 'install',

    template: Handlebars.compile(require('text!tmpl/install.hbs')),

    events: {
      'click .js-back': 'back',
      'click .js-new-user': 'newUser',
      'click .js-returning-user': 'returningUser',
      'click .js-set-key': 'setKey',
      'click .js-set-name': 'setName',
      'click .js-set-avatar': 'setAvatar'
    },

    initialize: function() {
      this.model.on('change:step', this.step, this);
    },

    newUser: function() {
      this.model.set('step', '2a');
    },

    returningUser: function() {
      this.model.set('step', '2b');
    },

    setName: function() {
      var name = this.$('.js-name').val();

      if (name) {
        this.model.set('name', name).set('step', '3');
      }
    },

    setKey: function() {
      var key = this.$('.js-key').val();

      if (key) {
        this.model.set('key', key).set('step', '3');
      }
    },

    setAvatar: function() {
      var avatar = this.$('.js-avatar').val();

      if (avatar) {
        this.model.set('avatar', avatar).set('step', '4');
      }
    },

    back: function() {
      var next;

      switch (this.model.get('step')) {
        case '2a':
        case '2b':
          this.model.unset('name', { silent: true });
          this.model.unset('key', { silent: true });
          next = '1';
          break;
        case '3':
          next = this.model.get('key') ? '2b' : '2a';
          break;
        case '4':
          next = '3';
      }

      this.model.set('step', next);
    },

    step: function() {
      this.$('.step').hide();
      this.$('.js-step-' + this.model.get('step')).show();

      if (this.model.get('step') === '4') {
        this.createUser();
      }
    },

    createUser: function() {
      logger.info('creating new user');
      logger.debug(this.model.attributes);

      var user = new User();
      user.set('name', this.model.get('name'));
      user.set('email', this.model.get('avatar'));

      user.save({}, {
        success: function() {
          vole.user = user;
          logger.info('successfully create new user');
          Backbone.history.navigate('/timeline', true);
        },
        error: function() {
          logger.error('unable to create user');
          this.model.set('step', '1');
        }
      });
    },

    spinner: function() {
      this.$('.spin').append(new Spin({
        width: 2
      }).spin().el);
    },

    render: function() {
      this.$el.html(this.template(this.model.attributes));
      this.spinner();
      this.step();
      return this;
    }

  });

});
