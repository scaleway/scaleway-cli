# Scaleway command-line bash completion
shopt -s extglob

# Requesting the servers/images in real time takes almost 1s which is long in a shell
# completion process - the following env variables (toggles) disable that behaviour
#GET_SCW_SRVS="no"
#GET_SCW_IMGS="no"

# Initialization
servers=""
images=""

# If you want to remove the inserted space after bash completion uncomment this line
# and comment the line just after it
#_nospace="-o nospace"
_nospace=""

_get_servers ()
{
    # Get a list of servers owned by the user.

    if [[ "${GET_SCW_SRVS}" =~ ^(yes|YES|0)$ ]]
    then
        servers=$(scw ps | awk '$0 !~ /SERVER/ {print $1}')
    fi
}

_get_images ()
{
    # Get a list of available images.

    if [[ "${GET_SCW_IMGS}" =~ ^(yes|YES|0)$ ]]
    then
        images=$(scw images | awk '$0 !~ /REPOSITORY/ {print $1}')
    fi
}

_scw_attach ()
{
    # Attach to a running server serial console.

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="--no-stdin
           -h
           --help"

    _get_servers

    COMPREPLY=( $(compgen -W "${words} ${servers}" -- ${cur}) )
    return 0
}

_scw_commit ()
{
    # Create a new snapshot from a server's volume.

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="-v
           --volume
           -h
           --help"

    _get_servers

    COMPREPLY=( $(compgen -W "${words} ${servers}" -- ${cur}) )
    return 0
}

_scw_cp ()
{
    # Copy files/folders from a PATH on the server to a HOSTDIR on the host
    # running the command. Use '-' to write the data as a tar file to STDOUT.

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="-g
           --gateway
           -p
           --port
           --user
           -h
           --help"

    _get_servers

    COMPREPLY=( $(compgen -W "${words} ${servers}" -- ${cur}) )
    return 0
}

_scw_create ()
{
    #Create a new server but do not start it.

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="--bootscript
           --commercial-type
           -e
           --env
           --ip-address
           --ipv6
           --name
           --tmp-ssh-key
           -v
           --volume
           -h
           --help"

    _get_images

    COMPREPLY=( $(compgen -W "${words} ${images}" -- ${cur}) )
    return 0
}

_scw_events ()
{
    # Get real time events from the API.

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="-h
           --help"

    COMPREPLY=( $(compgen -W "${words}" -- ${cur}) )
    return 0
}

_scw_exec ()
{
    # Run a command on a running server.

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="--A
           --g
           --gateway
           -p
           --port
           -T
           --timeout
           --user
           -w
           --wait
           -h
           --help"

    _get_servers

    COMPREPLY=( $(compgen -W "${words} ${servers}" -- ${cur}) )
    return 0
}

_scw_history ()
{
    # Show the history of an image.

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="--arch
           --no-trunc
           -q
           --quiet
           -h
           --help"

    _get_images

    COMPREPLY=( $(compgen -W "${words} ${images}" -- ${cur}) )
    return 0
}

_scw_images ()
{
    # List images.

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="-a
           --all
           -f
           --filter
           --no-trunc
           -q
           --quiet
           -h
           --help"

    COMPREPLY=( $(compgen -W "${words}" -- ${cur}) )
    return 0
}

_scw_info ()
{
    # Display system-wide information.

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="-h
           --help"

    COMPREPLY=( $(compgen -W "${words}" -- ${cur}) )
    return 0
}

_scw_inspect ()
{
    # Return low-level information on a server, image, snapshot, volume or bootscript.

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="--arch
           -b
           --browser
           -f
           --format
           -h
           --help"

    COMPREPLY=( $(compgen -W "${words}" -- ${cur}) )
    return 0
}

_scw_kill ()
{
    # Kill a running server.

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="-g
           --gateway
           -p
           --port
           -u
           --user
           -h
           --help"

    _get_servers

    COMPREPLY=( $(compgen -W "${words} ${servers}" -- ${cur}) )
    return 0
}

_scw_login ()
{
    # Generates a configuration file in '/home/$USER/.scwrc'
    # containing credentials used to interact with the Scaleway API. This
    # configuration file is automatically used by the 'scw' commands.
    #
    # You can get your credentials on https://cloud.scaleway.com/#/credentials

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="-o
           --organization
           -s
           --skip-ssh-key
           -t
           --token
           -h
           --help"

    COMPREPLY=( $(compgen -W "${words}" -- ${cur}) )
    return 0
}

_scw_logout ()
{
    # Log out from the Scaleway API.

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="-h
           --help"

    COMPREPLY=( $(compgen -W "${words}" -- ${cur}) )
    return 0
}

_scw_logs ()
{
    # Fetch the logs of a server.

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="-g
           --gateway
           -p
           --port
           --user
           -h
           --help"

    _get_servers

    COMPREPLY=( $(compgen -W "${words} ${servers}" -- ${cur}) )
    return 0
}

_scw_port ()
{
    # List port mappings for the SERVER, or lookup the public-facing port that is NAT-ed to the PRIVATE_PORT

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="-g
           --gateway
           -p
           --port
           --user
           -h
           --help"

    _get_servers

    COMPREPLY=( $(compgen -W "${words} ${servers}" -- ${cur}) )
    return 0
}

_scw_products ()
{
    # Display products PRODUCT information. At the moment only `servers` is supported.

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="servers
           -s
           --short
           -h
           --help"

    COMPREPLY=( $(compgen -W "${words}" -- ${cur}) )
    return 0
}

_scw_ps ()
{
    # List servers. By default, only running servers are displayed.

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="-a
           --all
           -f
           --filter
           -l
           --latest
           -n
           --no-trunc
           -q
           --quit
           -h
           --help"

    COMPREPLY=( $(compgen -W "${words}" -- ${cur}) )
    return 0
}

_scw_rename ()
{
    # Rename a server.

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="-h
           --help"

    _get_servers

    COMPREPLY=( $(compgen -W "${words} ${servers}" -- ${cur}) )
    return 0
}

_scw_restart ()
{
    # Restart a running server.

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="-T
           --timeout
           -w
           --wait
           -h
           --help"

    _get_servers

    COMPREPLY=( $(compgen -W "${words} ${servers}" -- ${cur}) )
    return 0
}

_scw_rm ()
{
    # Remove one or more servers.

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="-f
           --force
           -h
           --help"

    _get_servers

    COMPREPLY=( $(compgen -W "${words} ${servers}" -- ${cur}) )
    return 0
}

_scw_rmi ()
{
    # Remove one or more image(s)/volume(s)/snapshot(s)

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="-h
           --help"

    COMPREPLY=( $(compgen -W "${words}" -- ${cur}) )
    return 0
}

_scw_run ()
{
    # Run a command in a new server.

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="-a
           --attach
           --bootscript
           --commercial-type
           -d
           --detach
           -e
           --env
           -g
           --gateway
           --ip-address
           --ipv6
           --name
           -p
           --port
           --rm
           --show-boot
           -T
           --timeout
           --tmp-ssh-key
           -u
           --userdata
           --user
           -v
           --volume
           -h
           --help"

    _get_images

    COMPREPLY=( $(compgen -W "${words} ${images}" -- ${cur}) )
    return 0
}

_scw_search ()
{
    # Search the Scaleway Hub for images.

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="--no-trunc
           -h
           --help"

    COMPREPLY=( $(compgen -W "${words}" -- ${cur}) )
    return 0
}

_scw_start ()
{
    # Start a stopped server.

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="--set-state
           -T
           --timeout
           -w
           --wait
           -h
           --help"

    _get_servers

    COMPREPLY=( $(compgen -W "${words} ${servers}" -- ${cur}) )
    return 0
}

_scw_stop ()
{
    # Stop a running server.

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="-t
           --terminate
           -w
           --wait
           -h
           --help"

    _get_servers

    COMPREPLY=( $(compgen -W "${words} ${servers}" -- ${cur}) )
    return 0
}

_scw_tag ()
{
    # Tag a snapshot into an image.

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="--arch
           --bootscript
           -h
           --help"

    COMPREPLY=( $(compgen -W "${words}" -- ${cur}) )
    return 0
}

_scw_top ()
{
    # Lookup the running processes of a server.

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="-g
           --gateway
           -p
           --port
           -h
           --help"

    _get_servers

    COMPREPLY=( $(compgen -W "${words} ${servers}" -- ${cur}) )
    return 0
}

_scw_version ()
{
    # Show the version information.

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="-h
           --help"

    COMPREPLY=( $(compgen -W "${words}" -- ${cur}) )
    return 0
}

_scw_wait ()
{
    # Block until a server stops.

    local cur prev grps words=() #split=false

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    COMPREPLY=()
    words="-h
           --help"

    _get_servers

    COMPREPLY=( $(compgen -W "${words} ${servers}" -- ${cur}) )
    return 0
}

_scw ()
{
    local cur prev grps words=() #split=false

    COMPREPLY=()
    options="-D
             --debug
             -V
             --verbose
             -q
             --quiet
             --sensitive
             -v
             --version
             --region"
    words="attach
           commit
           cp
           create
           events
           exec
           history
           images
           info
           inspect
           kill
           login
           logout
           logs
           port
           products
           ps
           rename
           restart
           rm
           rmi
           run
           search
           start
           stop
           tag
           top
           version
           wait"

    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    # trigger var is used to check if bash completes when there is only options on the line
    trigger="no"

    case ${prev} in
        scw)
            COMPREPLY=( $(compgen -W "${words} ${options}" -- ${cur}) )
        ;;

        @($(echo ${words} | tr ' ' '|')) )
            _scw_${prev}
        ;;

        @($(echo ${options} | tr ' ' '|')) )
            COMPREPLY=( $(compgen -W "${words} ${options}" -- ${cur}) )
        ;;

        *)
            for word in ${COMP_WORDS[*]}
            do
                if [[ ${word} =~ ^($(echo ${words} | tr ' ' '|')) ]]
                then
                    trigger="yes"
                    _scw_${word}
                fi
            done
            if [[ ${trigger} != "yes" ]]
            then
                COMPREPLY=( $(compgen -W "${words} ${options}" -- ${cur}) )
            fi
        ;;
    esac

    return 0
} &&
complete -F _scw ${_nospace} -o noquote -o filenames scw
