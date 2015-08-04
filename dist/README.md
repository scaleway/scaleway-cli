# Run in Docker

You can run scaleway-cli in a sandboxed way using Docker.

**warning: caching is disabled**

```console
$ docker run -it --rm --volume=$HOME/.scwrc:/root/.scwrc scaleway/cli ps
```
