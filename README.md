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

Note that to avoid an egg and chicken problem, `nuv` itself is build with his ancestor, `task`.

- Build it with just `task build`.
- Run tests with `task test`.

# Documentation

## Where `nuv` looks for tasks

Nuv is an enhanced task runner. Tasks are described by [task](https://taskfile.dev) taskfiles.
Nuv is able either to run existing tasks or download them from github.

When you run `nuv [<args>...]` it will first look for its `nuv` root.

Once the `nuv` root has been found, it will download related tools and finally will execute tasks.

### How `nuv` locate the `nuv` root

The `nuv` root is a folder with two files in it: `Nuvfile` (an yaml taskfile) and `Nuvtools` (a json file describing the tools to download). 

The first step is to locate the task. The algorithm is the following.

First it will look in the current folder if there is a `Nuvfile`. If there is, it will also look for `Nuvtools`. If it is not there, it will go up of one level looking for a directory with `Nuvfile` and `Nuvtools`, and selects it as the `nuv` root.

If there is not a `Nuvfile` it will look for a folder called `nuvtasks` with both a `Nuvfile` and `Nuvtools` in it and selects it as the `nuv` root.

If the preceding tests fails, it will try to download it from GitHub, from a branch in a github repo.

The repo defaults to  `https://github.com/nuvolaris/nuvfiles` but it can be changes setting the enviroment variable `NUV_TASKS`. The branch to use is wired in the build and it is changed at build time. 

It will download or update `nuvtasks` and store the in `~/nuv/tasks`. Note it will update the `nuvtasks` if the repo is older that 24hours.

### Download tools

The `Nuvtools` is a json file in the format: 
```
{
  "<tool>": "<url>"
}
```

where the `<tool>` is the name of a multiplatform binary, and `<url>`is the url with some replacement strings. 

TODO: provide more details on the `Nuvtools` format

## How `nuv` execute tasks

It will then look to the command line parameters parameters `nuv <arg1> <arg2> <arg3>` and will consider them directory names. The list can be empty. 

If there is a directory name  `<arg1>` it will change to that directory. If there is then a subdirectory `<arg2>` it will change to that and so on until it finds a argument that is not a directory name. 

If the last argument is a directory name, will look for a `help.txt`. If it is there, it will show this. It is is not there it will execute a `task -l` showing tasks with description. 

If it finds an argument not corresponding to a directory, it will  consider it a task to execute.

It will xecute the `help.txt` as a  [`docopt`](http://docopt.org/) to parse the remaining arguments as parametes. It will return a series of `<key>=<value>` that will be fed to `task`.

So to make an example a command like `nuv setup kubernetes install --context=k3s` will look in the folder `setup/kubernetes` in the `nuv root` if it is there, select `install` as task to execute and parse the `--context=k3s`. It is equivalent to invoke `cd setup/kubernetes ; task install -- context=k3s`.

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
