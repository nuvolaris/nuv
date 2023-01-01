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
cd "$(dirname $0)"/tests

ARGS=""
function parse_args {
   HELP="$1.help"
   shift
   if [[ "$1" = "--help" ]]
   then  if [[ -e "$HELP" ]]
         then cat "$HELP"
         else echo "no help for $1"
         fi
         return 1 
   else  ARGS="$@"
         return 0
   fi
}

execute() {
   cmd="$1"
   shift
   help="$cmd.help"
   #echo helps is $help
   if test -e "$help"
   then echo ---
        echo "$cmd" "$@"
        echo ---
        docopts -h "$(cat $help)" : "$cmd" "$@"
        echo ---
        #echo "parse then execute" task "$cmd" -- "$@"
   else echo "execute" task "$cmd" -- "$@"
   fi
}

while true
do
   if [[ -z "$1" ]]
   then  task
         break
   elif [[ "$1" = -* ]] 
   then  if parse_args default "$@"
         then execute default "$@"
         fi 
         break 
   elif [[ -d "$1" ]]
   then  cd "$1"
         shift
         continue
   else  A="$1"
         shift
         if parse_args "$A" "$@"
         then  execute "$A" "$ARGS"
         fi
         break
   fi
done
