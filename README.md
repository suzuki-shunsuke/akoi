# akoi

[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/suzuki-shunsuke/akoi)
[![Go Report Card](https://goreportcard.com/badge/github.com/suzuki-shunsuke/akoi)](https://goreportcard.com/report/github.com/suzuki-shunsuke/akoi)
[![GitHub last commit](https://img.shields.io/github/last-commit/suzuki-shunsuke/akoi.svg)](https://github.com/suzuki-shunsuke/akoi)
[![GitHub tag](https://img.shields.io/github/tag/suzuki-shunsuke/akoi.svg)](https://github.com/suzuki-shunsuke/akoi/releases)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/suzuki-shunsuke/akoi/master/LICENSE)

binary installer

## ⚠️ Deprecated: Please use aqua instead of akoi :warning:

https://aquaproj.github.io/

---

* [Overview](#overview)
* [Features](#features)
* [Getting Started](#getting-started)
* [Install](#install)
* [Usage](#usage)
* [Configuration](#configuration)
* [Output Format](#output-format)
  * [ansible](#ansible)
* [Tips](#tips)
  * [Use akoi at ansible](#use-akoi-at-ansible)
* [Release Notes](https://github.com/suzuki-shunsuke/akoi/releases)
* [License](#license)

## Overview

`akoi` is a binary installer.
`akoi` installs binaries according to the configuration file.

`akoi`'s task is simple.

1. Download archives
2. Unarchive downloaded archives
3. Install binaries
4. Create symbolic links to binaries

`akoi` manages the binary's version with symbolic link.
For example,

```
/usr/local/bin/consul-1.2.1
/usr/local/bin/consul -> consul-1.2.1
```

## Features

* efficiency
  * install packages in parallel
  * support parallel download with Accept-Ranges
* declarative and idempotence
* [small dependencies and easy to install (written in Go)](#install)
* [work good with ansible's shell module](#use-akoi-at-ansible)

## Demo

<p align="center">
  <img src="https://cdn.jsdelivr.net/gh/suzuki-shunsuke/artifact@master/akoi/demo.gif">
</p>

## Getting Started

[Install akoi](#install).
Generate the configuration file.

```
$ akoi init --dest akoi.yml
```

Edit the `akoi.yml`.

In this section, install [consul](https://www.consul.io/) as example.

```yaml
---
# akoi - binary version control system
# https://github.com/suzuki-shunsuke/akoi
bin_path: /tmp/bin/{{.Name}}-{{.Version}}
link_path: /tmp/bin/{{.Name}}
packages:
  consul:
    url: "https://releases.hashicorp.com/consul/{{.Version}}/consul_{{.Version}}_{{.OS}}_{{.Arch}}.zip"
    version: 1.2.1
    files:
    - name: consul
      archive: consul
```

Run the `akoi install` command to install binaries.

```
$ akoi install -c akoi.yml
downloading consul: https://releases.hashicorp.com/consul/1.2.1/consul_1.2.1_darwin_amd64.zip
unarchive consul
create directory /tmp/bin
install /tmp/bin/consul-1.2.1
create link /tmp/bin/consul -> /tmp/bin/consul-1.2.1
```

Binaries are installed at the `bin` directory.

Check bin directory.

```
$ ls /tmp/bin
consul -> consul-1.2.1
consul-1.2.1
```

Edit `akoi.yml` to change the consul version to 1.2.0 and run `akoi install` again.

```
$ akoi install -c akoi.yml
downloading consul: https://releases.hashicorp.com/consul/1.2.0/consul_1.2.0_darwin_amd64.zip
unarchive consul
install /tmp/bin/consul-1.2.0
remove link /tmp/bin/consul -> /tmp/bin/consul-1.2.1
create link /tmp/bin/consul -> /tmp/bin/consul-1.2.0
```

Run `akoi install` again. `akoi` does nothing. `akoi` doesn't download files wastefully.

```
$ akoi install -c akoi.yml # output nothing
```

Edit `akoi.yml` to change the consul version to 1.2.1 and run `akoi install` again.
The consul 1.2.1 has already been installed so `akoi` doesn't download the archive wastefully.
`akoi` only recreates the symbolic link.

```
$ akoi install -c akoi.yml
remove link /tmp/bin/consul -> /tmp/bin/consul-1.2.0
create link /tmp/bin/consul -> /tmp/bin/consul-1.2.1
```

## Install

akoi is written with Golang and binary is distributed at [release page](https://github.com/suzuki-shunsuke/akoi/releases), so installation is easy and no dependency is needed.

If you want to build yourself, run the following command. 

```
$ go get -u github.com/suzuki-shunsuke/akoi
```

If you want to install `akoi` with ansible, please consider to use the ansible role [suzuki-shunsuke.akoi](https://galaxy.ansible.com/suzuki-shunsuke/akoi).

Check whether akoi is installed.

```
$ akoi -v
akoi version 0.3.1
```

## Usage

Please run `akoi help` or `akoi help <command>`.

```
$ akoi help
$ akoi help [init|install]
```

## Template Engine

The following fields in the configuration file are parsed by [Go's text/template](https://golang.org/pkg/text/template/).

* bin_path
* link_path
* archive
* url

And in template, [sprig functions](http://masterminds.github.io/sprig/) can be used.

## Configuration

```yaml
---
# bin_path and link_path are absolute path or relative path from the configuration file parent directory.
# binary install path
bin_path: $HOME/bin/{{.Name}}-{{.Version}}
# the symbolic link to the binary
link_path: $HOME/bin/{{.Name}}
# number of download partitions with Accept-Ranges
# Note that this parameter is ignored if server doesn't support Accept-Ranges.
# If this parameter is not set or is less equal than 0, the value of `runtime.NumCPU()` is used.
# If this parameter is 1, Accept-Ranges isn't used.
num_of_dl_partitions: 4
max_parallel_download_count: 5
http_request_timeout: 300
packages:
  consul: # package name
    # akoi downloads a file from this url and unarchive it according to the base file name.
    # akoi uses https://github.com/mholt/archiver to unarchive the file.
    # .OS and .Arch are got by runtime.GOOS and runtime.GOARCH .
    # https://golang.org/pkg/runtime/#pkg-constants
    url: "https://releases.hashicorp.com/consul/{{.Version}}/consul_{{.Version}}_{{.OS}}_{{.Arch}}.zip"
    # package version
    version: 1.2.1
    # archive file's type. This is optional and by default this is decided by url's path.
    # akoi uses mholt/archiver and the value included in archiver.SupportedFormat is supported.
    # https://github.com/mholt/archiver
    # If downloaded file is not archived, set archive_type to "unarchived".
    # And the file type is compressed but not archived (not ".tar.gz"), set archive_type to "Gzip".
    archive_type: Zip
    # binary install path
    bin_path: /usr/local/bin/{{.Name}}-{{.Version}}
    # the symbolic link to the binary
    link_path: /usr/local/bin/{{.Name}}
    num_of_dl_partitions: 1
    http_request_timeout: 120
    # files included in the downloaded file
    files:
    - name: consul
      # when unarchive the file to the temporary directory,
      # the relative path from the directory to the file.
      archive: consul
      # file's mode. This is optional and default value is 0755.
      mode: 0644
      # binary install path
      bin_path: /opt/bin/{{.Name}}-{{.Version}}
      # the symbolic link to the binary
      link_path: /opt/bin/{{.Name}}
```

## Environment variables

### AKOI_CONFIG_PATH

Configuration file path. The precedence is

1. command line option
2. AKOI_CONFIG_PATH
3. finds the file .akoi.yml from the current directory to the root directory recursively
4. /etc/akoi/akoi.yml

## Output Format

### ansible

*Note that this specification is unstable.*

```json
{
  "msg": "",
  "changed": true,
  "failed": false,
  "packages": {
    "consul": {
      "name": "consul",
      "changed": true,
      "failed": false,
      "error": "",
      "version": "1.2.1",
      "url": "https://releases.hashicorp.com/consul/1.2.1/consul_1.2.1_darwin_amd64.zip",
      "files": {
        "consul": {
          "name": "consul",
          "error": "",
          "changed": true,
          "migrated": false,
          "installed": true,
          "mode_changed": false,
          "file_removed": false,
          "dir_created": false,
          "link": "consul",
          "entity": "consul"
        }
      }
    }
  }
}
```

## Tips

### Use akoi at ansible

If you want to install binaries with ansible, run `akoi install` command with `--format` option in [ansible's shell](https://docs.ansible.com/ansible/latest/modules/shell_module.html) module.
When `--format` option is set `akoi install` outputs the result as json.

```
$ akoi install -f ansible | jq '.'
{
  "msg": "",
  "changed": false,
  "packages": {
    "consul": {
      "error": "",
      "changed": false,
      "state": "",
      "files": [
        {
          "error": "",
          "changed": false,
          "state": "",
          "name": "consul",
          "link": "",
          "entity": ""
        }
      ],
      "version": "1.2.1",
      "url": ""
    }
  }
}
```

So you can check whether the task's result by passing the output.

```yaml
tasks:
- name: install consul
  shell: "/usr/local/bin/akoi install -f ansible 2>&1"
  register: result
  changed_when: (result.stdout|from_json)["changed"]
```

## Change Log

Please see [Releases](https://github.com/suzuki-shunsuke/akoi/releases).

## License

[MIT](LICENSE)
