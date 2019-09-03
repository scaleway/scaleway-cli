# anonuuid

:wrench: anonymize UUIDs

![anonuuid Logo](https://raw.githubusercontent.com/moul/anonuuid/master/assets/anonuuid.png)

[![CircleCI](https://circleci.com/gh/moul/anonuuid.svg?style=shield)](https://circleci.com/gh/moul/anonuuid)
[![GoDoc](https://godoc.org/moul.io/anonuuid?status.svg)](https://godoc.org/moul.io/anonuuid)
[![License](https://img.shields.io/github/license/moul/anonuuid.svg)](https://github.com/moul/anonuuid/blob/master/LICENSE)
[![GitHub release](https://img.shields.io/github/release/moul/anonuuid.svg)](https://github.com/moul/anonuuid/releases)
[![Go Report Card](https://goreportcard.com/badge/moul.io/anonuuid)](https://goreportcard.com/report/moul.io/anonuuid)
[![CodeFactor](https://www.codefactor.io/repository/github/moul/anonuuid/badge)](https://www.codefactor.io/repository/github/moul/anonuuid)
[![codecov](https://codecov.io/gh/moul/anonuuid/branch/master/graph/badge.svg)](https://codecov.io/gh/moul/anonuuid)
[![Docker Metrics](https://images.microbadger.com/badges/image/moul/anonuuid.svg)](https://microbadger.com/images/moul/anonuuid)
[![Sourcegraph](https://sourcegraph.com/github.com/moul/anonuuid/-/badge.svg)](https://sourcegraph.com/github.com/moul/anonuuid?badge)
[![Made by Manfred Touron](https://img.shields.io/badge/made%20by-Manfred%20Touron-blue.svg?style=flat)](https://manfred.life/)


**anonuuid** anonymize an input string by replacing all UUIDs by an anonymized
new one.

The fake UUIDs are cached, so if `anonuuid` encounter the same real UUIDs multiple
times, the translation will be the same.

## Usage

```console
$ anonuuid --help
NAME:
   anonuuid - Anonymize UUIDs outputs

USAGE:
   anonuuid [global options] command [command options] [arguments...]

VERSION:
   1.0.0-dev

AUTHOR(S):
   Manfred Touron <https://moul.io/anonuuid>

COMMANDS:
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --hexspeak		Generate hexspeak style fake UUIDs
   --random, -r		Generate random fake UUIDs
   --keep-beginning	Keep first part of the UUID unchanged
   --keep-end		Keep last part of the UUID unchanged
   --prefix, -p 	Prefix generated UUIDs
   --suffix 		Suffix generated UUIDs
   --help, -h		show help
   --version, -v	print the version
   ```

## Example

Replace all UUIDs and cache the correspondance.

```command
$ anonuuid git:(master) ✗ cat <<EOF | anonuuid
VOLUMES_0_SERVER_ID=15573749-c89d-41dd-a655-16e79bed52e0
VOLUMES_0_SERVER_NAME=hello
VOLUMES_0_ID=c245c3cb-3336-4567-ada1-70cb1fe4eefe
VOLUMES_0_SIZE=50000000000
ORGANIZATION=fe1e54e8-d69d-4f7c-a9f1-42069e03da31
TEST=15573749-c89d-41dd-a655-16e79bed52e0
EOF
VOLUMES_0_SERVER_ID=00000000-0000-0000-0000-000000000000
VOLUMES_0_SERVER_NAME=hello
VOLUMES_0_ID=11111111-1111-1111-1111-111111111111
VOLUMES_0_SIZE=50000000000
ORGANIZATION=22222222-2222-2222-2222-222222222222
TEST=00000000-0000-0000-0000-000000000000
```

---

Inline

```command
$ echo 'VOLUMES_0_SERVER_ID=15573749-c89d-41dd-a655-16e79bed52e0 VOLUMES_0_SERVER_NAME=bitrig1 VOLUMES_0_ID=c245c3cb-3336-4567-ada1-70cb1fe4eefe VOLUMES_0_SIZE=50000000000 ORGANIZATION=fe1e54e8-d69d-4f7c-a9f1-42069e03da31 TEST=15573749-c89d-41dd-a655-16e79bed52e0' | ./anonuuid
VOLUMES_0_SERVER_ID=00000000-0000-0000-0000-000000000000 VOLUMES_0_SERVER_NAME=bitrig1 VOLUMES_0_ID=11111111-1111-1111-1111-111111111111 VOLUMES_0_SIZE=50000000000 ORGANIZATION=22222222-2222-2222-2222-222222222222 TEST=00000000-0000-0000-0000-000000000000
```

---

```command
$ curl -s https://api.pathwar.net/achievements\?max_results\=2 | anonuuid | jq .
{
  "_items": [
    {
      "_updated": "Thu, 30 Apr 2015 13:00:58 GMT",
      "description": "You",
      "_links": {
        "self": {
          "href": "achievements/00000000-0000-0000-0000-000000000000",
          "title": "achievement"
        }
      },
      "_created": "Thu, 30 Apr 2015 13:00:58 GMT",
      "_id": "00000000-0000-0000-0000-000000000000",
      "_etag": "b1e9f850accfcb952c58384db41d89728890a69f",
      "name": "finish-20-levels"
    },
    {
      "_updated": "Thu, 30 Apr 2015 13:01:07 GMT",
      "description": "You",
      "_links": {
        "self": {
          "href": "achievements/11111111-1111-1111-1111-111111111111",
          "title": "achievement"
        }
      },
      "_created": "Thu, 30 Apr 2015 13:01:07 GMT",
      "_id": "11111111-1111-1111-1111-111111111111",
      "_etag": "c346f5e1c4f7658f2dfc4124efa87aba909a9821",
      "name": "buy-30-levels"
    }
  ],
  "_links": {
    "self": {
      "href": "achievements?max_results=2",
      "title": "achievements"
    },
    "last": {
      "href": "achievements?max_results=2&page=23",
      "title": "last page"
    },
    "parent": {
      "href": "/",
      "title": "home"
    },
    "next": {
      "href": "achievements?max_results=2&page=2",
      "title": "next page"
    }
  },
  "_meta": {
    "max_results": 2,
    "total": 46,
    "page": 1
  }
}
```

## Install

### Using go

- `go get moul.io/anonuuid/cmd/anonuuid`

### Using brew

- `brew install moul/moul/anonuuid`

### Download release

https://github.com/moul/anonuuid/releases

## License

© 2015-2019 [Manfred Touron](https://manfred.life) -
[Apache-2.0 License](https://github.com/moul/anonuuid/blob/master/LICENSE)
