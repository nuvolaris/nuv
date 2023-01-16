#!/bin/bash
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#
cd "$(dirname $0)/olaris"

TMP=/tmp/args$$
DEBUG=""
if test "$1" = "-d"
then shift ; DEBUG=1
fi

while test -e "$1/Taskfile.yml"
do cd "$1" ; shift
done

if test -z "$1" || test "$1" = "help" 
then task $1
else 
    if cat help.txt | docopts -h - : $cmd "$@" >/tmp/args$$
    then 
        test -n "$DEBUG" && cat $TMP | xargs echo task "$1"
        cat $TMP | xargs task "$1"
    else cat $TMP | sh
    fi
    rm -f $TMP
fi

