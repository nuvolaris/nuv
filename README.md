<!--
  ~ Licensed to the Apache Software Foundation (ASF) under one
  ~ or more contributor license agreements.  See the NOTICE file
  ~ distributed with this work for additional information
  ~ regarding copyright ownership.  The ASF licenses this file
  ~ to you under the Apache License, Version 2.0 (the
  ~ "License"); you may not use this file except in compliance
  ~ with the License.  You may obtain a copy of the License at
  ~
  ~   http://www.apache.org/licenses/LICENSE-2.0
  ~
  ~ Unless required by applicable law or agreed to in writing,
  ~ software distributed under the License is distributed on an
  ~ "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
  ~ KIND, either express or implied.  See the License for the
  ~ specific language governing permissions and limitations
  ~ under the License.
  ~
-->
**WARNING: this is still work in progress**

The code may not build, there can be errors that will destroy your hard disk and so on.

Also documentation and code may not be aligned, as we first write documentation and then write code, so you may read documentation of code that does not yet exist!

You have been warned!

# `nuv`, the next generation

Nuv is the nuvolaris all-mighty build tool.

This is a rewrite of the current build tools to make it super powerful.

It is basically the [task](https://taskfile.dev) tool enanced to support:

- a bunch of embedded commands (check tooks) including `wsk` 
- the ability to download other tools on the file
- a predefined set of tasks downloaded from github
- a way to create a hierarchy of taskfiles 
- documentation for tasks powered by [docopt](http://docopt.org/)

Note that to avoid an egg and chicken problem, `nuv` itself is built with his ancestor, `task`.

- Build it with just `task build`.
- Run tests with `task test`.

# Documentation

## Environment variables

The following environment variables allows to ovverride certain defaults.

- `NUV_ROOT` is the  folder where `nuv` looks for its tasks. It is not defined, it follows the algorithm below to find it.
- `NUV_BIN` is the  folder where `nuv` looks binaries (external command line tools). It is not defined, it defaults to the same directory where  `nuv` is located.
- `NUV_REPO` is the github repo where `nuv` downloads its tasks. It it is not defined, it defaults to `https://github.com/nuvolaris/olaris`
- `NUV_BRANCH` is the branch where `nuv` looks for its tasks. The branch to use is defined at build time and it is the base version (without the patch level). For example, if `nuv` is `0.3.0-morpheus` the branch to use will be `0.3-morpheus`

## Where `nuv` looks for tasks

Nuv is an enhanced task runner. Tasks are described by [task](https://taskfile.dev) taskfiles.
Nuv is able either to run existing tasks or download them from github.

When you run `nuv [<args>...]` it will first look for its `nuv` root.

Once the `nuv` root has been found, it will download related tools and finally will execute tasks.

### How `nuv` locate the `nuv` root

The `nuv` root is a folder with two files in it: `nuvfile.yml` (an yaml taskfile) and `nuvtools.json` (a json file describing the tools to download). 

The first step is to locate the root folder. The algorithm is the following.

If the environment variable `NUV_ROOT` is defined, it will look there first, and will check if there are the two files.

Then it will look in the current folder if there is a `nuvfile.yml`. If there is, it will also look for `nuvtools.json`. If it is not there, it will go up of one level looking for a directory with `nuvfile.yml` and `nuvtools.json`, and selects it as the `nuv` root.

If there is not a `nuvfile.yml` it will look for a folder called `olaris` with both a `nuvfile.yml` and `nuvtools.json` in it and will select it as the `nuv` root.

If the preceding tests fails, it will try to download it from GitHub, from a branch in a github repo. 

The repo to use is defined by the environment variable `NUV_REPO`, and defaults if it is missing to `https://github.com/nuvolaris/olaris`

The branch to use it defined at build time. It can be overriden with the enviroment variable `NUV_BRANCH`.

If it is missing, it will clone the repo and store the in `~/.nuv/olaris`. If it is there and the `nuvtools.yml` is older than 24 hours, it will update it. Otherwise it will stop the search there.

## Command Line Tools

Nuv requires some binary command line tools to work with ("bins").

They are expected to be in the folder pointed by the environment variable `NUV_BIN`. 

If this environment variable is not defined, it defaults to the same folder where `nuv` itself is located. The `NUV_BIN` folder is then added to the beginning of the `PATH` before executing anything else.

When nuv update it will check the tools are there and of the right version running `nuv -t nuvtools.yml check`, and issue a warning if there is a version mismatch.

Updating tools can be executed manually with `nuv -t nuvtools.yaml update`, but generally you need to be  root or administrator.

Nuv is however normally distributed with an installer that includes all the tools for the various operating systems (linux, windows, osx).

## How `nuv` execute tasks

It will then look to the command line parameters parameters `nuv <arg1> <arg2> <arg3>` and will consider them directory names. The list can be empty. 

If there is a directory name  `<arg1>` it will change to that directory. If there is then a subdirectory `<arg2>` it will change to that and so on until it finds a argument that is not a directory name. 

If the last argument is a directory name, will look for a `nuvopts.txt`. If it is there, it will show this. It is is not there it will execute a `task -t nuvfile.yml -l` showing tasks with description. 

If it finds an argument not corresponding to a directory, it will consider it a task to execute, 

If there is not a `nuvopts.txt`,  it will execute as a task, passing the other arguments (equivalent to `task -t nuvfile.yml <arg> -- <the-other-args>`). 

If there is a `nuvopts.txt`, it will interpret it as a  [`docopt`](http://docopt.org/) to parse the remaining arguments as parameters. The result of parsing is a sequence of `<key>=<value>` that will be fed to `task`. So it is equivalent to invoking `task -t nuvfile.yml <arg> <key>=<value> <key>=<value>...`

### Example

A command like `nuv setup kubernetes install --context=k3s` will look in the folder `setup/kubernetes` in the `nuv root` if it is there, select `install` as task to execute and parse the `--context=k3s`. It is equivalent to invoke `cd setup/kubernetes ; task install -- context=k3s`.

Note that also this will also use the downloaded tools and the embedded commands of `nuv`.

## Embedded tools

Currently task embeds the following tools, and you can invoke them directly prefixing them with `-`: (`nuv -task`, `nuv -basename` etc). Use `nuv -help` to list them.

- [task](https://taskfile.dev) the Task build tool
- [wsk](https://github.com/apache/openwhisk-cli) the OpenWhisk cli 
- [ht](https://github.com/nojima/httpie-go) an httpie tool in Golan

Basic unix like tools (`nuv -<tool> -help for details`):

- basename
- cat
- cp
- dirname
- grep
- gunzip
- gzip
- head
- ls
- mv
- pwd
- rm
- sleep
- tail
- tar
- tee
- touch
- tr
- unzip
- wc
- which
- zip







