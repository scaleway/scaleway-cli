var Api = require('onlinelabs'),
    Table = require('cli-table'),
    _ = require('lodash');


module.exports.panic = function(msg) {
  console.error(msg);
  process.exit(-1);
};


module.exports.notImplementedAction = function() {
  console.error("onlinelabs: Not implemented");
};


module.exports.truncateRow = function(row, limits) {
  for (idx in row) {
    if (limits[idx] != -1) {
      row[idx] = row[idx].toString().substring(0, limits[idx]);
    }
  }
  return row;
};


module.exports.defaultConfigPath = function() {
  var home = process.env[(
    process.platform == 'win32' ?
      'USERPROFILE' :
      'HOME'
  )];
  return home + '/.onlinelabsrc';
};


module.exports.newTable = function(options) {
  options = options || {};
  options.chars = options.chars || {
    'top': '', 'top-mid': '', 'top-left': '', 'top-right': '',
    'bottom': '', 'bottom-mid': '', 'bottom-left': '', 'bottom-right': '',
    'left': '', 'left-mid': '', 'mid': '', 'mid-mid': '',
    'right': '', 'right-mid': '', 'middle': ' '
  };
  options.style = options.style || {
    // 'padding-left': 0, 'padding-right': 0
  };
  return new Table(options);
};


module.exports.wordify = function(str) {
  return str
    .replace(/[^a-zA-Z0-9-]/g, '_')
    .replace(/__+/g, '_')
    .replace(/^_/, '')
    .replace(/_$/, '');
};


module.exports.newApi = function(options) {
  // var config = _.defaults(options || {}, rc);
  var config = module.exports.rc();
  return new Api(config);
};


module.exports.collect = function(val, memo) {
  memo.push(val);
  return memo;
};


module.exports.rc = function() {
  return require('onlinelabs/node_modules/rc')('onlinelabs', {
    api_endpoint: 'https://api.cloud.online.net/',
    organization: null,
    token: null
  });
};
