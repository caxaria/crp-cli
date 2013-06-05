require('colors');
var read = require('read');
var adminClient = require('crp-admin-client')();

module.exports = signup;

module.exports.usage =
function usage(name, args) {
  args.
    alias('i', 'invitation').
    usage('crowdprocess signup [<email> [<password>] [-i <invitation>]]');
};

function signup(args) {
  username(function(username) {
    password(function(password) {
      doSignup(username, password, args.invitation, function(err) {
        if (err) {
          console.error(err.message.red);
          return;
        }
        console.log('Signed up successfully'.green);
      });
    });
  });

  function username(cb) {
    var username = args._[0];
    if (username) return cb(username);

    read({
      prompt: 'Email:'
    }, function(err, username) {
      if (err) throw err;
      cb(username);
    });
  }

  function password(cb) {
    var password = args._[1];
    if (password) return cb(password);

    read({
      prompt: 'Password:',
      silent: true
    }, function(err, password) {
      if (err) throw err;
      cb(password);
    });
  }
}

function doSignup(username, password, invitation, cb) {
  adminClient.account.create(username, password, invitation, cb);
}