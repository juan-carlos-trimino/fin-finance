package misc

import (
  "os/user"
  "runtime"
  "strings"
)

//Is the current user running as root?
func IsRoot() (bool, error) {
  current, err := user.Current()
  if err != nil {
    return false, err
  }
  return strings.EqualFold(current.Username, "root"), nil
}

func GetUsername() (string, error) {
  current, err := user.Current()
  if err != nil {
    return "", err
  }
  return current.Username, nil
}

func GetOS() string {
  return runtime.GOOS
}

/***
When you create a new file or directory, it is assigned default permissions. There are two things
that affect the default permissions. The first is whether you are creating a regular file or a
directory; the second is the current umask.

If you create a new regular file, the OS assigns it the default permissions (octal) 0666
(-rw-rw-rw-). If you create a new directory, the OS assigns it the default permissions (octal) 0777
(drwxrwxrwx). However, the shell session will also set a umask to further restrict the permissions
that are initially set. This is an octal bitmask used to clear the permissions of new files and
directories created by a process. If a bit is set in the umask, then the corresponding permission
is cleared on new files. The default mask for non-root users is 0002 and for root users is 0022.
Hence, when you create a file with the default settings for a root user, the final file permissions
will be 0644 (rw-r--r--), and for a directory, the final permissions will be 755 (rwxr-xr-x).

Three permission categories apply: read, write, and execute. The following table explains how these
permissions affect access to files and directories.

Effects of Permissions on Files and Directories
---------------------------------------------------------------------------------------------------
Permission  Effect on files                     Effective on directories
---------------------------------------------------------------------------------------------------
r (read)    File contents can be read.          Contents of the directory (the file names) can be
                                                listed.
w (write)   File contents can be changed.       Any file in the directory can be created, deleted,
                                                or renamed.
x (execute) Files can be executed as commands.  The directory can become the current working
                                                directory.

The umask command without arguments will display the current value of the shell's umask:
$ umask

To change the value the shell's umask:
$ umask 0660
You can omit any leading zeros in the umask.

The system's default umask values for Bash shell users are defined in the /etc/profile and
/etc/bashrc files. Users can override the system defaults in the .bash_profile and .bashrc files
in their home directories.

To get the permissions for a file:
$ ls -l file_name.txt

And to get the permissions for a directory:
$ ls -ld directory_name
***/
