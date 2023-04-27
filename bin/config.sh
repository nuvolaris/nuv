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

NUV_TMP=~/.nuv/tmp

if test -z "$1"
then cat $NUV_TMP/env 
     exit 0
fi

for I in "$@"
do
    IFS='=' read -r K V <<<"$I"
    if test -z "$V"
    then rm $NUV_TMP/$K.env
    else echo $I >$NUV_TMP/$K.env
    fi
done
cat $NUV_TMP/*.env >$NUV_TMP/env
