#!/bin/bash
printf "Usage:\n"
printf "./inspect_pod_node.sh namespace [label_key=label_value]\n"
# ${parameter:+word}
# If parameter is null or unset, nothing is substituted, otherwise the expansion of word is
# substituted. The value of parameter is not used.
NAMESPACE=${1:+"-n ${1}"}
# When shift is used without an argument, it shifts the positional parameters by one position. This
# means the value of $2 becomes $1, $3 becomes $2, and so on. The original value of $1 is
# discarded.
# You can provide a positive integer as an argument to shift (e.g., shift 2). This will shift the
# parameters by that specified number of positions. For instance, shift 2 would make $3 the new $1,
# $4 the new $2, and so forth. The value of the positive integer must be between zero and the
# number of positional parameters ($#), inclusive.
# The special variable $#, which represents the total number of positional parameters, is also
# updated by shift to reflect the new count after the shift operation.
shift
APP_SELECTOR=${1:+"--selector=${1}"}
shift
# Save the original IFS value.
declare -a ARR=($(kubectl get pods ${APP_SELECTOR} ${NAMESPACE} \
                  --field-selector status.phase!=Pending \
                  -o jsonpath='{range .items[*]}{.metadata.name}{","}{.spec.nodeName}{"\n"}{end}'))
ORIGINAL_IFS="$IFS"
# printf "%q\n" "$IFS"
# In Bash scripting, IFS stands for Internal Field Separator (or sometimes Input Field Separator).
# It is a special shell variable that defines the characters used to delimit words or fields when
# Bash performs word splitting.
# By default, IFS is set to a space, a tab, and a newline character (<space><tab><newline>). This
# means that when Bash reads input or expands variables, it will split the string into separate
# "words" or "fields" wherever it encounters any of these characters.
# Set IFS to a comma.
IFS=","
printf "The columns in the output represent the hostname, cloud zone, and Pod name.\n\n"
# Iterate through each of the pods and output Pod name and node name.
for ITEM in "${ARR[@]}"
do
  # For each Pod/Node line item, obtain Cloud Zone for Node.
  POD_NODE=($ITEM)
  kubectl get node ${POD_NODE[1]} \
    -o jsonpath='{.metadata.name}{"\t"}{.metadata.labels.topology\.kubernetes\.io\/zone}{"\t"}{"'${POD_NODE[0]}'"}{"\n"}'
done
# Restore IFS to its default value.
IFS="$ORIGINAL_IFS"
