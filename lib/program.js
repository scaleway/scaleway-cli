var Q = require('q'),
    _ = require('lodash'),
    async = require('async'),
    child_process = require('child_process'),
    debug = require('debug')('onlinelabs-cli:program'),
    filesize = require('filesize'),
    filesizeParser = require('filesize-parser'),
    fs = require('fs'),
    jsonPath = require('JSONPath'),
    moment = require('moment'),
    program = require('commander'),
    termJsCli = require('../node_modules/term.js-cli'),
    utils = require('./utils');


program
  .version(utils.getVersion('..'))
  .option('--api-endpoint <url>', 'set the API endpoint')
  .option('-D, --debug', 'enable debug mode');


program
  .command('attach <server>')
  .description('attach (serial console) to a running server')
  .action(function(server, options) {
    var client = utils.newApi(options);
    utils.searchEntity({
      input: server,
      filters: {
        _type: 'servers'
      }
    }, function(err, entity) {
      var ttyUrl = 'https://tty.cloud.online.net?server_id=' + entity._id
            + '&type=serial&auth_token=' + client.config.token;
      termJsCli(ttyUrl, function(err) {
        if (err) { utils.panic(err); }
      });
    });
  });


program
  .command('build <path>')
  .description('build an image from a file')
  .action(utils.notImplementedAction);


program
  .command('commit <server> [name]')
  .description("create a new snapshot from a server's volume")
  .option('-a, --author <author>',
          'author (e.g., "Georges Abitbol <georges@most-class.world>")')
  .option('--name <name>', 'assign a name to the snapshot', 'noname')
  .option('-p, --pause', 'pause server during commit')
  .option('-v, --volume <slot>', 'volume slot')
  .action(function(server, name, options) {
    var client = utils.newApi(options);
    var volumeIdx = options.volume || 0;

    if (options.author || options.pause) {
      utils.panic("onlinelabs: Not implemented option");
    }

    utils.searchEntity({
      input: server,
      filters: {
        _type: 'servers'
      }
    }, function(err, entity) {
      client.get('/servers/' + entity._id)
        .then(function(res) {
          name = options.name || (res.body.server.name + '-snapshot');
          client.post('/snapshots', {
            volume_id: res.body.server.volumes[volumeIdx.toString()].id,
            organization: res.body.server.organization,
            name: name
          })
            .then(function(res) {
              utils.saveEntity(res.body.snapshot, 'snapshots');
              console.log(res.body.snapshot.id);
            })
            .catch(utils.panic);
        })
        .catch(utils.panic);
    });
  });


program
  .command('cp <server:path> <path>')
  .description("copy files/folders from a server's filesystem to the host path")
  .action(utils.notImplementedAction);


program
  .command('create <image>')
  .description('create a new server but do not start it')
  .option('--name <name>', 'assign a name to the server', 'noname')
  .option('--bootscript <bootscript>', 'assign a bootscript')
  .option('-v, --volume <size>', 'attach additional volume', utils.collect, [])
  .option('-e, --env <environments>',
          'provide metadata tags passed to initrd (i.e., boot=rescue,INITRD_DEBUG=1)',
          utils.collect, [])
  .action(function(image, options) {
    var client = utils.newApi(options);

    var createServer = function(options) {
      var data = {
        organization: client.config.organization,
        name: options.name
      };
      if (options.image) {
        data.image = options.image;
      }
      if (options.bootscript) {
        data.bootscript = options.bootscript;
      }
      if (options.root_volume) {
        data.volumes = data.volumes || {};
        data.volumes['0'] = options.root_volume;
      }
      if (options.volumes && options.volumes.length) {
        data.volumes = data.volumes || {};
        options.volumes.forEach(function(volume, idx) {
          data.volumes[(idx + 1).toString()] = volume;
        });
      }
      if (options.env) {
        data.tags = options.env;
      }

      client.post('/servers', data)
        .then(function (res) {
          utils.saveEntity(res.body.server, 'servers');
          console.log(res.body.server.id);
        })
        .catch(utils.panic);
    };

    // Create volumes
    Q.all(_.map(options.volume, function(volume) {
      return client.post('/volumes', {
        organization: client.config.organization,
        size: parseInt(filesizeParser(volume, {base: 10})),
        name: volume,
        volume_type: 'l_ssd'
      });
    })).then(
      function(results) {
        var volumes = _.pluck(_.pluck(results, 'body'), 'volume');
        options.volumes = _.pluck(volumes, 'id');
        _.forEach(volumes, function(volume) {
          utils.saveEntity(volume, 'volumes');
        });
        // Resolve bootscript
        utils.searchEntity({ input: options.bootscript, filters: { _type: 'bootscripts' } }, function(err, bootscriptEntity) {
          if (options.bootscript && err) { utils.panic(err); }
          options.bootscript = bootscriptEntity.id;

          // Resolve image
          utils.searchEntity({ input: image, filters: { _type: 'images' } }, function(err, imageEntity) {
            if (err) {
              var size = filesizeParser(image, {base: 10});
              if (!size) {
                utils.panic(err);
              }
              client.post('/volumes', {
                organization: client.config.organization,
                size: parseInt(size),
                name: image,
                volume_type: 'l_ssd'
              }).then(function(results) {
                options.root_volume = results.body.volume.id;
                createServer(options);
              }).catch(utils.panic);
            } else {
              options.image = imageEntity.id;
              createServer(options);
            }
          });
        });
      });
  });


program
  .command('events')
  .description('get real time events from the API')
  .option('-f, --filter <filters>',
          'provide filter values. valid filters: (i.e., status=pending)',
          utils.collect, [])
  .option('--since <timestamp>', 'show all events created since timestamp')
  .option('--until <timestamp>', 'stream events until this timestamp')
  .action(function(options) {
    var client = utils.newApi(options);
    client.get('/tasks')
      .then(function(res) {
        _.forEach(res.body.tasks, function(task) {
          console.log(task.started_at + ' ' +
                      task.href_from + ': ' +
                      task.description + ' ('+
                      task.status + ' ' +
                      task.progress + ') ' +
                      task.terminated_at);
        });
      })
      .catch(utils.panic);
  });


program
  .command('exec <server> <command> [args...]')
  .description('run a command in a running server')
  .option('-d, --detach', 'detached mode: run command in the background')
  .option('-i, --interactive', 'keep STDIN open even if not attached')
  .option('-t, --tty', 'allocate a pseudo-TTY')
  .option('-k, --insecure', 'disable SSH strict host key checking')
  .action(function(server, command, commandArgs, options) {
    var client = utils.newApi(options);
    utils.searchEntity({
      input: server,
      filters: {
        _type: 'servers'
      }
    }, function(err, entity) {
      if (err) { utils.panic(err); }
      client.get('/servers/' + entity._id)
        .then(function(res) {
          if (!res.body.server.public_ip) {
            utils.panic('Server ' + res.body.server.id + ' is not running');
          }
          var ip = res.body.server.public_ip.address;

          var args = [];

          if (!debug.enabled) {
            args.push('-q');
          }

          args = args.concat.apply(args, [
            '-l', 'root',
            ip, '-t', '--', command], commandArgs);

          if (options.insecure) {
            args = [].concat.apply([
              '-o', 'UserKnownHostsFile=/dev/null',
              '-o', 'StrictHostKeyChecking=no'
            ], args);
          }

          debug('spawn: ssh ' + args.join(' '));
          var spawn = child_process.spawn(
            'ssh',
            args,
            { stdio: 'inherit' }
          );
          spawn.on('close', function(code) {
            process.exit(code);
          });
        })
        .catch(utils.panic);
    });
  });


program
  .command('export <server>')
  .description('stream the contents of a server as a tar archive')
  .action(utils.notImplementedAction);


program
  .command('history <image>')
  .description('show the history of an image')
  .option('--no-trunc', "don't truncate output")
  .option('-q, --quiet', 'only display numeric IDs')
  .action(function(image, options) {
    var client = utils.newApi(options);
    utils.searchEntity({
      input: image,
      filters: {
        _type: 'images'
      }
    }, function(err, entity) {
      if (err) {
        utils.panic(err);
      }
      client.get('/images/' + entity._id)
        .then(function(res) {
          if (options.quiet) {
            console.log(res.body.image.id);
          } else {
            var table = utils.newTable({
              head: [
                'IMAGE', 'CREATED', 'CREATED BY', 'SIZE'
              ]
            });

            var image = res.body.image;
            var row = [
              image.id,
              moment(image.creation_date).fromNow(),
              image.root_volume.name,
              filesize(image.root_volume.size, {base: 10})
            ];
            if (options.trunc) {
              utils.truncateRow(row, [8, 25, 25, 25]);
            }
            table.push(row);
            console.log(table.toString());
          }
        })
        .catch(utils.panic);
    });
  });


program
  .command('images')
  .description('list images')
  .option('-a, --all', 'show all images')
  .option('-f, --filter <filters>',
          "provide filter values. (i.e., 'public=true')", utils.collect, [])
  .option('--no-trunc', "don't truncate output")
  .option('-q, --quiet', 'only display numeric IDs')
  .action(function(options) {
    var client = utils.newApi(options);
    var promises = [];

    var query = '/images?';
    if (options.filter.length) {
      utils.panic("onlinelabs: Not implemented option");
    }
    promises.push(client.get(query));

    if (options.all) {
      promises.push(client.get('/snapshots'));
      promises.push(client.get('/bootscripts'));
    }

    Q.all(promises).then(
      function(results) {
        var entries = _.reduce(
          _.pluck(results, 'body'),
          function(entries, group) {
            return entries.concat.apply(
              entries,
              _.reduce(
                group,
                function(aggreg, n, key) {
                  utils.saveEntities(n, key);
                  return aggreg.concat.apply(
                    aggreg,
                    _.map(n, function(entry) {
                      entry._type = key;
                      return entry;
                    })
                  );
                }, [])
            );
          }, []);

        if (options.quiet) {
          _.forEach(
            _.sortByOrder(entries, ['creation_date'], [false]),
            function(entry) {
              console.log(entry.id);
            });
        } else {
          var table = utils.newTable({
            head: [
              'REPOSITORY', 'TAG', 'IMAGE ID', 'CREATED', 'VIRTUAL SIZE'
            ]
          });

          _.forEach(_.sortByOrder(
            entries,
            ['creation_date'],
            [false]
          ), function(entry) {
            var repository, tag, imageId, created, virtualSize;
            switch (entry._type) {
            case 'snapshots':
              var snapshot = entry;
              repository = utils.wordify(snapshot.name);
              tag = '<none>';
              imageId = snapshot.id;
              created = moment(snapshot.creation_date).fromNow();
              virtualSize = filesize(snapshot.size, {base: 10});
              break;
            case 'images':
              var image = entry;
              repository = utils.wordify(image.name);
              if (!image.public) {
                repository = 'user/' + utils.wordify(image.name);
              }
              tag = 'latest';
              imageId = image.id;
              created = moment(image.creation_date).fromNow();
              virtualSize = filesize(image.root_volume.size, {base: 10});
              break;
            case 'bootscripts':
              var bootscript = entry;
              repository = utils.wordify(bootscript.title);
              tag = 'bootscript';
              imageId = bootscript.id;
              created = 'n/a';
              virtualSize = 'n/a';
              break;
            }
            var row = [
              repository, tag, imageId, created, virtualSize
            ];
            if (options.trunc) {
              utils.truncateRow(row, [40, 25, 8, 25, 25]);
            }
            table.push(row);
          });
          console.log(table.toString());
        }

      }, utils.panic);
  });


program
  .command('import')
  .description('create a new filesystem image from the contents of a tarball')
  .action(utils.notImplementedAction);


program
  .command('info')
  .description('display system-wide information')
  .action(function() {
    var rc = utils.rc();
    console.log('Organization: ' + rc.organization);
    console.log('Token: ' + utils.anonymizeUUID(rc.token));
    console.log('API Endpoint: ' + rc.api_endpoint);
    console.log('RC file: ' + rc.config);
    console.log('CLI path: ' + process.argv[1]);
    console.log('User: ' + process.env.USER);
    utils.db.count({}, function(err, count) {
      if (!err) {
        console.log('Cached entities: '+ count);
      }
    });
  });


program
  .command('inspect <items...>')
  .description('return low-level information on a server or image')
  .option('-f, --format <format>', 'format the output using the given template')
  .action(function(items, options) {
    var client = utils.newApi(options);
    var promises = [];

    var once = function(item, cb) {
      return [
        client.get('/servers/' + item._id),
        client.get('/images/' + item._id),
        client.get('/volumes/' + item._id),
        client.get('/bootscripts/' + item._id)
      ];
    };

    utils.searchEntities({
      inputs: items,
      filters: {}
    }, function(err, entities) {
      if (err) {
        utils.panic(err);
      }
      promises = promises.concat.apply(promises, entities.map(once));

      Q.allSettled(promises).then(
        function(results) {
          var entries = _.filter(
            _.pluck(
              _.pluck(
                results,
                'value'
              ),
              'body'
            ),
            function(entry) {
              return entry !== undefined;
            });

          if (options.format) {
            _.map(entries, function(entry) {
              var parsed = jsonPath.eval(entry, '$' + options.format);
              if (typeof(parsed) === 'object' && parsed.length === 1) {
                console.log(parsed[0]);
              } else {
                console.log(parsed);
              }
            });
          } else {
            console.log(JSON.stringify(entries, null, 2));
          }
        }, utils.panic);
    });
  });


program
  .command('kill <server>')
  .description('kill a running server')
  .option('-s, --signal <signal>', 'Signal to send to the server', 'KILL')
  .action(function(server, options) {
    if (options.signal !== 'KILL') {
      utils.panic("onlinelabs: Not implemented option");
    }
    var client = utils.newApi(options);

    utils.searchEntity({
      input: server,
      filters: {
        _type: 'servers'
      }
    }, function(err, entity) {
      if (err) {
        utils.panic(err);
      }
      client.get('/servers/' + entity._id)
        .then(function(res) {
          if (!res.body.server.public_ip) {
            utils.panic('Server ' + res.body.server.id + ' is not running');
          }
          var ip = res.body.server.public_ip.address;

          utils.sshExec(ip, 'halt', {}, function(statusCode) {
            if (statusCode === 0) {
              console.log(server);
            }
            process.exit(statusCode);
          });
        })
        .catch(utils.panic);
    });
  });


program
  .command('load')
  .description('load an image from a tar archive')
  .action(utils.notImplementedAction);


program
  .command('login')
  .description('login to the API')
  .option('--organization <uuid>', 'set the organization')
  .option('--token <token>', 'token')
  .action(function(options) {
    var client = utils.newApi(options);
    var newConfig = _.cloneDeep(client.config);
    delete newConfig._;
    delete newConfig.configs;
    delete newConfig.config;
    var filepath = utils.defaultConfigPath();
    fs.writeFile(
      filepath,
      JSON.stringify(newConfig, null, 2),
      function (err) {
        if (err) {
          utils.panic(err);
        }
        console.log('configuration written to ' + filepath);
      });
  });


program
  .command('logout')
  .description('log out from the API')
  .action(function() {
    var filepath = utils.defaultConfigPath();
    fs.unlink(
      filepath,
      function (err) {
        if (err) {
          utils.panic(err);
        }
        console.log('removed ' + filepath + ' configuration file');
      });
  });


program
  .command('logs <server>')
  .description('fetch the logs of a server')
  .action(utils.notImplementedAction);


program
  .command('port')
  .description('list port security for the server')
  .action(utils.notImplementedAction);


program
  .command('pause')
  .description('pause all processes within a server')
  .action(utils.notImplementedAction);


program
  .command('ps')
  .description('list servers')
  .option('-a, --all',
          'show all servers. only running servers are shown by default')
  .option('--before <server>', 'show only server created before server, ' +
          'include non-running ones')
  .option('-f, --filter <filters>', 'provide filter values. valid filters: ' +
          'status=(starting|running|stopping|stopped)', utils.collect, [])
  .option('-l, --latest',
          'show only the latest created server, include non-running ones')
  .option('-n <n>', 'show n last created servers, include non-running ones.',
          parseInt)
  .option('--no-trunc', "don't truncate output")
  .option('-q, --quiet', 'only display numeric IDs')
  .option('-s, --size', 'display total file sizes')
  .option('--since <server>',
          'show only servers created since server, include non-running ones')
  .action(function(options) {
    var client = utils.newApi(options);
    var query = '/servers?';

    if (options.before || options.filter.length || options.size ||
        options.since) {
      utils.panic("onlinelabs: Not implemented option");
    }

    if (!options.all)   { query += 'state=running&'; }
    if (options.latest) { query += 'per_page=1&'; }
    if (options.n)      { query += 'per_page=' + options.n + '&'; }

    client.get(query)
      .then(function(res) {
        if (options.all) {
          utils.saveEntities(res.body.servers, 'servers');
        } else {
          // FIXME: saveEntity
        }
        if (options.quiet) {
          _.forEach(
            _.sortByOrder(res.body.servers, ['creation_date'], [false]),
            function(server) {
              console.log(server.id);
            });
        } else {
          var table = utils.newTable({
            head: [
              'SERVER ID', 'IMAGE', 'COMMAND', 'CREATED', 'STATUS', 'PORTS',
              'NAME'
            ]
          });

          _.forEach(_.sortByOrder(
            res.body.servers,
            ['creation_date'],
            [false]), function(server) {
            var row = [
              server.id,
              (server.image ? utils.wordify(server.image.root_volume.name) : ''),
              '',
              moment(server.creation_date).fromNow(),
              server.state,
              '',
              utils.wordify(server.name)
            ];
            if (options.trunc) {
              utils.truncateRow(row, [8, 25, 25, 25, 25, 25, -1]);
            }
            table.push(row);
          });
          console.log(table.toString());
        }
      })
      .catch(utils.panic);
  });


program
  .command('pull <image>')
  .description('pull an image or a repository')
  .action(utils.notImplementedAction);


program
  .command('push <image>')
  .description('push an image or a repository')
  .action(utils.notImplementedAction);


program
  .command('rename <server>')
  .description('rename an existing server')
  .action(utils.notImplementedAction);


program
  .command('restart <server>')
  .description('restart a running server')
  .option('-t, --time <second>', 'number of seconds to try to stop for ' +
          'before killing the server. once killed it will be restarted.')
  .action(function(server, options) {
    var client = utils.newApi(options);

    if (options.time) {
      utils.panic("onlinelabs: Not implemented option");
    }

    utils.searchEntity({
      input: server,
      filters: {
        _type: 'servers'
      }
    }, function(err, entity) {
      client.post('/servers/' + entity._id + '/action', {
        action: 'reboot'
      })
        .then(function() {
          console.log(server);
        })
        .catch(function (err) {
          if (err.error.message !== 'server is being stopped or rebooted') {
            utils.panic(err);
          }
        });
    });
  });


program
  .command('rm <server>')
  .description('remove one or more servers')
  .action(function(server, options) {
    var client = utils.newApi(options);
    utils.searchEntity({
      input: server,
      filters: {
        _type: 'servers'
      }
    }, function(err, entity) {
      client.delete('/servers/' + entity._id)
        .then(function(res) {
          if (res.statusCode !== 204) {
            utils.panic(res);
          }
        }).catch(utils.panic);
    });
  });


program
  .command('rmi <image>')
  .description('remove one or more images')
  .action(function(image, options) {
    var client = utils.newApi(options);
    utils.searchEntity({
      input: image,
      filters: {
        _type: 'images'
      }
    }, function(err, entity) {
      client.delete('/images/' + entity._id)
        .then(function(res) {
          if (res.statusCode !== 204) {
            utils.panic(res);
          }
        }).catch(utils.panic);
    });
  });


program
  .command('run <image>')
  .description('run a command in a new server')
  .action(utils.notImplementedAction);


program
  .command('save <image>')
  .description('save an image to a tar archive')
  .action(utils.notImplementedAction);


program
  .command('search <keyword>')
  .description('search for an image on the Hub')
  .action(utils.notImplementedAction);


program
  .command('start <server>')
  .description('start a stopped server')
  .option('-a, --attach', "attach server's STDOUT and STDERR and forward " +
          'all signals to the process')
  .option('-i, --interactive', "attach server's STDIN")
  .option('-s, --sync', 'synchronous start. wait for SSH to be ready')
  .action(function(server, options) {
    var client = utils.newApi(options);
    var serverId = server;

    if (options.attach || options.interactive) {
      utils.panic("onlinelabs: Not implemented option");
    }

    utils.searchEntity({
      input: serverId,
      filters: {
        _type: 'servers'
      }
    }, function(err, entity) {
      client.post('/servers/' + entity._id + '/action', {
        action: 'poweron'
      })
        .then(function() {
          console.log(serverId);
          if (options.sync) {
            utils.waitForServerState(client, entity._id, 'running', function(err, server) {
              utils.waitPortOpen(server.public_ip.address, 22, function(err) {
                if (err) { utils.panic(err); }
              });
              debug('server state is running');
            });
          }
        })
        .catch(function (err) {
          if (err.error.message !== 'server should be stopped') {
            utils.panic(err);
          }
        });
    });
  });


program
  .command('stop <server>')
  .description('stop a running server')
  .option('-t, --terminate', 'stop and trash a server and its volumes')
  .action(function(server, options) {
    var client = utils.newApi(options);

    var data = {
      action: 'poweroff'
    };
    if (options.terminate) {
      data.action = 'terminate';
    }

    utils.searchEntity({
      input: server,
      filters: {
        _type: 'servers'
      }
    }, function(err, entity) {
      client.post('/servers/' + entity._id + '/action', data)
        .then(function() {
          console.log(server);
        })
        .catch(function (err) {
          if (!_.includes([
            'server is being stopped or rebooted',
            'server should be running'
          ], err.error.message)) {
            utils.panic(err);
          }
        });
    });
  });


program
  .command('tag <snapshot> <tag-name>')
  .description('tag an image into a repository')
  .action(function(snapshot, tagName, options) {
    var client = utils.newApi(options);

    utils.searchEntity({
      input: snapshot,
      filters: {
        _type: 'snapshots'
      }
    }, function(err, entity) {
      client.get('/snapshots/' + entity._id)
        .then(function(res) {
          client.post('/images', {
            root_volume: res.body.snapshot.id,
            organization: res.body.snapshot.organization,
            name: tagName,
            arch: 'arm'
          })
            .then(function(res) {
              utils.saveEntity(res.body.image, 'images');
              console.log(res.body.image.id);
            })
            .catch(utils.panic);
        })
        .catch(utils.panic);
    });
  });


program
  .command('top <server>')
  .description('lookup the running processes of a server')
  .action(utils.notImplementedAction);


program
  .command('unpause <server>')
  .description('unpause a paused server')
  .action(utils.notImplementedAction);


program._events.version = null;
program
  .command('version')
  .description('show the version information')
  .action(function() {
    console.log('Client version: ' + utils.getVersion('..'));
    console.log('Client API version: ' + utils.getVersion('onlinelabs'));
    console.log('Node.js version (client): ' + process.version);
    console.log('OS/Arch (client): ' + process.platform + '/' + process.arch);
    // FIXME: add information about server
  });



program
  .command('wait <server>')
  .description('block until a server stops')
  .action(function(server, options) {
    var client = utils.newApi(options);

    utils.searchEntity({
      input: server,
      filters: {
        _type: 'servers'
      }
    }, function(err, entity) {
      utils.waitForServerState(client, entity._id, 'stopped', function(err) {
        if (err) {
          utils.panic(err);
        }
        console.log(0);
      });
    });
  });


module.exports = program;
