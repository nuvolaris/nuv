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
         if [[ -n "$ARGS" ]]
         then echo "[args $ARGS]"
         fi
         return 0
   fi
}

while true
do
   if [[ -z "$1" ]]
   then  task
         break
   elif [[ "$1" = -* ]] 
   then  if parse_args default "$@"
         then task default -- "$ARGS"
         fi 
         break 
   elif [[ -d "$1" ]]
   then  cd "$1"
         shift
         continue
   else  A="$1"
         shift
         if parse_args "$A" "$@"
         then task "$A" -- "$ARGS"
         fi
         break
   fi
done
