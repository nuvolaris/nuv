#!/bin/bash

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
