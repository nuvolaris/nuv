#!/bin/bash
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

