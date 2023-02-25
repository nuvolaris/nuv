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

When you run `nuv` it will first look in current directory for a `Nuvfile.yml`. If it is not there, it will assume it has to download its tasks from github and will execute or update the `~/.nuv/olaris` by cloning or updating `https://github.com/nuvolaris/olaris` and change to that directory to look for tasks, and restart looking.

If it finds the `Nuvfile.yml`, it will also search for `Nuvtools.json`, a json file that lists the tools that must be downloaded (`kind`, `kubectl` etc). If there is not a `Nuvtools.json` it will go up one level and restart looking.

When finally it found a `Nuvfile.yml` with `Nuvtools.json`, it downloads all the tools and starts execution with that directory as the starting poing.

It will then look to the command line parameters parameters `nuv <arg1> <arg2> <arg3>` and will consider them directory names. The list can be empty. If there is a directory name  `<arg1>` it will change to that directory. If there is then a subdirectory `<arg2>` it will change to that and so on until it finds a argument that does not look a directory name. 

If it is not a directory name, and there are no other arguments, will look for a "help.txt". If it is there, it will show this. It is is not there it will execute a `task -l` showing tasks with description. 

If it finds a task with no directory, if there is the `help.txt` it will feed the remaining arguments to `docopt` to parse the parameters and finally will invoke the task with all the parameters interpreted according the  docopt specifications. Otherwise will just execute the task and use the remaining arguments as task parameters.

Phew!

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
