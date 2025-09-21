#!/bin/sh
#

# Notes
# (1) When attempting to run a script directly without specifying the interpreter (e.g.,
#     /bin/bash script.sh) or if the script's shebang line (#!/bin/sh) is incorrect or the
#     interpreter is missing in the container, the script will not run, and the deployment will
#     fail.
#
# (2) An exec probe's output is not sent to the pod's stdout and will not appear in kubectl logs.
#     The output is internal to the kubelet and is only visible when a probe fails. If you need to
#     debug the probe's script, you must redirect its output to a file inside the container.
#
#     Why probe output is not visible
#     The exec probe is executed directly by the kubelet on the node, not by the container's shell
#     or main process. The kubelet uses the command's exit code (0 for success, non-zero for
#     failure) to determine the container's health. The standard output is captured by the kubelet
#     but is not forwarded to the pod's logging stream.
#
#     How to debug a probe script
#     To see the output of your script for debugging, modify your probes to redirect the output to
#     a file that you can later inspect.
#
#     Check the probe logs from a running pod
#     After deploying the pod, you can use kubectl exec to view the contents of the log file
#     written by the probe script.
#     $ kubectl exec -n finances <pod-name> -- cat /wsf_data_dir/health-logs/probes.log
#
#     Alternative: View failing probe events
#     For a failed probe, the script's output will be recorded in the pod's events. You can view
#     this with kubectl describe.
#     $ kubectl describe pod -n finances <pod-name>
#     Under the Events section, you will see a Readiness/Liveness probe failed event.
#
# (3) When a K8s HTTP probe targets the wrong port, it will return an HTTP status code of 0. This
#     status code indicates that the client (the kubelet, in this case) did not receive any
#     response from the server within the specified timeout period, rather than receiving an actual
#     HTTP response.
#     In summary, an HTTP status code of 0 from a probe directed to the wrong port signifies a
#     fundamental failure to establish a connection or receive any data from the target, rather
#     than an application-level error like a 404 or 500 status.

MAX_SIZE_BYTES=102400  # 100KB in bytes
N_LAST_LINES=100  # Keep the last N lines
# The "-c %s" option specifically extracts and displays the file size in bytes.
FILE_SIZE=$(stat -c %s "$2")
if [ "$FILE_SIZE" -gt "$MAX_SIZE_BYTES" ];
then
  # Using 'tail' to keep the last N_LAST_LINES lines of the file.
  tail -n "$N_LAST_LINES" "$2" > "$2_tmp" && mv "$2_tmp" "$2"
  printf "\n\nFile size ($FILE_SIZE) exceeds the limit ($MAX_SIZE_BYTES). Keeping " >> "$2"
  printf "the last $N_LAST_LINES lines of the file.\n" >> "$2"
  printf "File truncated. New size: %d bytes.\n" "$(stat -c %s "$2")" >> "$2"
fi
# Explanation of the command:
# -S: This option prints the server response headers.
# -q: Turn off wget's output.
# -O /dev/null: This redirects the standard output (e.g., download progress) to /dev/null,
#               preventing it from cluttering your terminal.
# -nv: Turn off verbose without being completely quiet (use -q for that), which means that error
#      messages and basic information still get printed.
# "$1": Replace this with the URL you want to check.
# 2>&1: This redirects standard error to standard output, ensuring that error messages (which might
#       contain the HTTP status code if there's a problem) are included in the pipeline.
# awk '/^  HTTP/{print $2}': This awk command filters the output. It searches for lines starting
#     with "  HTTP" (two spaces followed by "HTTP", which is typical in wget's server response
#     output) and then prints the second field on that line, which corresponds to the HTTP status
#     code.
# status_code=$(wget ...): This captures the output of the entire pipeline into the status_code
#                          variable.
status_code=$(wget -nv -S -q -O /dev/null "$1" 2>&1 | awk '/^  HTTP/{print $2}')
printf "\n%s\n" "$(date -u)" >> "$2"
if [[ "$1" == *"liveness"* ]];
then
  printf "Status code for liveness probe: %d\n" "$status_code" >> "$2"
else
  printf "Status code for readiness probe: %d\n" "$status_code" >> "$2"
fi
#
if [ "$status_code" -ne 200 ];
then
  exit 1  # Non-zero exit code for failure.
fi
exit 0  # Zero exit code for success.
