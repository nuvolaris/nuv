#!/bin/bash

cur="$1"
while test -n "$cur"
do 
   if [[ "$cur" =~ "^-.*" ]]
   then  echo "parse args"
         break
   elif test -d "$cur"
   then  cd "$cur"
         shift
         cur="$1"
   else  break
   fi
done
if test -z "$cur"
then task -- "$@"
elif [[ "$cur" = "--help" ]]
then echo "help!"
else echo "parse args: $@"
      task "$cur"
fi
