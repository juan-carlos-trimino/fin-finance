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
export APP_VERSION=1.0.0
echo "APP_VERSION=$(printenv APP_VERSION)"
export K8S_MANIFEST_CRD=true
echo "K8S_MANIFEST_CRD=$(printenv K8S_MANIFEST_CRD)"
