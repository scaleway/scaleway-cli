OnlineLabs CLI
==============

Interact with OnlineLabs API from the command line.


Usage
-----

Usage 100% inspired by Docker

    $ onlinelabs

      Usage: onlinelabs [options] [command]


      Commands:

        attach <server>                 attach (serial console) to a running server
        build <path>                    build an image from a file
        commit <server>                 create a new image from a server's changes
        cp <server:path> <path>         copy files/folders from a server's filesystem to the host path
        create [options] <image>        create a new server but do not start it
        events                          get real time events from the API
        exec <server> <command>         run a command in a running server
        export <server>                 stream the contents of a server as a tar archive
        history <image>                 show the history of an image
        images [options]                list images
        import                          create a new filesystem image from the contents of a tarball
        info                            display system-wide information
        inspect <item> [otherItems...]  return low-level information on a server or image
        kill                            kill a running server
        load                            load an image from a tar archive
        login [options]                 login to the API
        logout                          log out from the API
        logs <server>                   fetch the logs of a server
        port                            list port security for the server
        pause                           pause all processes within a server
        ps [options]                    list servers
        pull <image>                    pull an image or a repository
        push <image>                    push an image or a repository
        rename <server>                 rename an existing server
        restart <server>                restart a running server
        rm <server>                     remove one or more servers
        rmi <image>                     remove one or more images
        run <image>                     run a command in a new server
        save <image>                    save an image to a tar archive
        search <keyword>                search for an image on the Hub
        start <server>                  start a stopped server
        stop <server>                   stop a running server
        tag <image> <tag>               tag an image into a repository
        top <server>                    lookup the running processes of a server
        unpause <server>                unpause a paused server
        version                         show the version information
        wait <server>                   block until a server stops

      Options:

        -h, --help            output usage information
        -V, --version         output the version number
        --api-endpoint <url>  set the API endpoint
        -D, --debug           enable debug mode


Install
-------

1. Install `Node.js` and `npm`
2. Install `onlinelabs-cli`: `$ npm install -g onlinelabs-cli`
3. Setup token and organization: `$ onlinelabs login --token=XXXXX --organization=YYYYY`


License
-------

MIT
