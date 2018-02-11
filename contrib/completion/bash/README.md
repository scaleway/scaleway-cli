# Bash completion of Scaleway command-line

This script helps you use `scw` command-line easily by generating a bash completion everytime you hit the `tab` key.

## How it works

- Simply copy the script named `scw.bash` to `/etc/bash_completion.d` it will be sourced everytime you login or you simulate a login.
- You can also put it in your `HOMEDIR` and source it with your login shell option by adding the following line in your `~/.bash_profile`:
```
source ~/.scw.bash
```

## Tuning variables

For more flexibility the behaviour of the completion if controlled by the following vars:

- **GET_SCW_SRVS**: it could be an env variable in your shell set by the command `export GET_SCW_SRVS="yes"`. If this option is enabled the completion will call Scaleway API to fetch for your managed servers and displays them. It default to "no" because fetching servers names takes around 1s which is quite long for command-line completion.
- **GET_SCW_IMGS**: it could be an env variable in your shell set by the command `export GET_SCW_IMGS="yes"`. If this option is enabled the completion will call Scaleway API to fetch for the available images and displays them. It default to "no" because fetching images takes around 1s which is quite long for command-line completion.
- **nospace**: this is a built-in shell completion option when set to `-o nospace` allows you to remove the space add to word when the completion happens.

## TODO

- Tune every function completion for more parsing options by tracking all the sub-options and handle them.
