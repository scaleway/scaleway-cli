"use strict";

var chai = require('chai'),
    debug = require('debug')('tests'),
    program = require('..'),
    stdout = require('test-console').stdout,
    util = require('util');


// Initialize chai.should()
chai.should();


var inspect = function(name, obj) {
  debug(name, util.inspect(obj, {showHidden: false, depth: null, colors: true}));
};


var run = function(command) {
  var args = [];
  args = args.concat.apply(['node', 'onlinelabs'], command);
  var inspect = stdout.inspect();
  program.parse(args);
  inspect.restore();
  return inspect.output.join('');
};


suite("[program]", function() {
  test('info', function() {
    var output = run(['info']);
    (output).should.contain('User: ' + process.env['USER']);
  });
  test('version', function() {
    var output = run(['version']);
    (output).should.contain('Client version: ' + require('../package.json').version);
    (output).should.contain('Client API version: ' + require('onlinelabs/package.json').version);
  });
});
