# pstalk

Creates arbitrarily deep process trees that do nothing.

## Motivation

I was working on [pocket][pocket] and learning about process termination, reparenting, etc., and I wanted an easy way create parent/child process relationships to experiemnt with.

## Usage

Run `pstalk <n>` to create a process tree that's `<n>` deep.

## Process Management Observations

- Killing any given process directly should kill only that process, i.e. killing an ancestor does not kill any of the descendants. This is the default bahvior for both Linux and Windows when killing a process and `pstalk` doesn't perform any process handling on its own.

- Killing the CLI process with ctrl+C kills all remaining processes on Linux and Windows, even if the lineage has been broken by killing a process directly. This is because terminals create _process groups_ for new processes, and child processes are by default members of the same groups, which allows the terminals to kill the entire process group. See [`setpgid`][setpgid] (Linux), [CREATE_NEW_PROCESS_GROUP][create-new-process-group], and [`GenerateConsoleCtrlEvent`][generateconsolectrlevent] (Windows) for more info.

- On Linux, killing the process group directly will kill all remaining `pstalk` processes in the tree even if the lineage was broken by an intermediate process exiting or being killed and leaving a grandchild isolated from its grandparent. This is because each descendant is created as a member of the group and group membership doesn't rely on parentage.

- On Windows there's no accessible way to interact with a process group. Though [the docs][create-new-process-group] say that "The process identifier of the new process group is the same as the process identifier" when creating a new process group, calling [`GenerateConsoleCtrlEvent`][generateconsolectrlevent] with the the root process ID, its parent ID, or any of the descendant process IDs returns `false` and doesn't stop any of the processes. The closest thing to native process group management is using `taskkill /f /t /pid <pid>` which kills the process with pid=`<pid>` and any accessible descendants. Unlike process management on Linux, the `taskkill` can kill subtrees by specifying a non-root `<pid>`, and doesn't traverse across breaks in the lineage caused by intermediate processes existing or being killed.

[create-new-process-group]: https://learn.microsoft.com/en-us/windows/win32/procthread/process-creation-flags

[generateconsolectrlevent]: https://learn.microsoft.com/en-us/windows/console/generateconsolectrlevent

[pocket]: https://github.com/ttd2089/pocket

[setpgid]: https://man7.org/linux/man-pages/man2/setpgid.2.html
