#!/bin/bash
#
# The "Dot Space Script" Trick for Running Bash Scripts in the "Current Shell"
# In Bash, a dot (".") is a command, which is synonymous with the "source" command. It executes the
# content of a specified script in the current shell environment. This differs from running a
# script directly, which executes in a subshell. Therefore, when using ".", any changes made by the
# script, such as variables assignments or directory changes, will persist in the current shell. To
# persist the environment variables, run the script like so
# $ . ./ev_app.sh
# or
# $ source ./ev_app.sh
# Please note that "source" is a synonym for dot (".") in the Bash shell, but not in POSIX sh, so
# for maximum compatibility use the period.
#
usage() {
  echo "Usage: . ./set-ev.sh -v APP_VERSION -rp true|fase"
  echo -e "Usage: . ./set-ev.sh --av APP_VERSION --reverse_proxy true|fase\n"
  return
}

# Check for an even number of arguments.
if (( $# % 2 != 0 ))
then
  echo -e "\nError: Please provide an even number of arguments."
  usage
else
  echo "Elements in \$@: $@"
  arr=()
  for arg in "$@"
  do
    arr+=("$arg")
  done
  # echo "Elements in arr: ${arr[@]}"
  declare -i is_error=0
  declare -i av_missing=0
  declare -i rp_missing=0
  declare -i ndx=0
  declare -i size="${#arr[@]}"
  declare app_version=""
  declare reverse_proxy="false"
  for (( ndx = 0; ndx < size & is_error == 0; ))
  do
    flag=${arr[ndx]}
    # echo "Element at index $ndx: $flag"
    ndx=$(( ndx + 1 ))
    value=${arr[ndx]}
    # echo "Element at index $ndx: $value"
    ndx=$(( ndx + 1 ))
    case "$flag" in
      "-av" | "--app_version")
        app_version=$value
        av_missing=1
        ;;
      "-rp" | "--reverse_proxy")
        reverse_proxy=$value
        rp_missing=1
        ;;
      *)  #Default.
        echo -e "\nUnknown flag $flag"
        usage
        is_error=1
        ;;
    esac
  done
  echo -e "\n"
  if [ $is_error -eq 0 ]
  then
    if (( av_missing == 1 ))
    then
      export APP_VERSION=$app_version
      # echo "APP_VERSION=$(printenv APP_VERSION)"
    fi
    #
    if [ ! -z "$reverse_proxy" ] && [ "$reverse_proxy" != "true" ] && [ "$reverse_proxy" != "false" ]
    then
      echo -e "Valid values for the flag -rp/--reverse_proxy: 'true' or 'false'.\n"
    elif (( rp_missing == 1 ))
    then
      export K8S_MANIFEST_CRD=$reverse_proxy
      # echo "K8S_MANIFEST_CRD=$(printenv K8S_MANIFEST_CRD)"
    fi
    echo "Current Values:"
    echo "APP_VERSION=$(printenv APP_VERSION)"
    echo -e "K8S_MANIFEST_CRD=$(printenv K8S_MANIFEST_CRD)\n"
  fi
fi
