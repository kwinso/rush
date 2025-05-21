# Rush
Rush is a simple shell implementation in Go programming language. It is pretty limited in it's functionality, but it was a fun project to learn about how Linux works on a deeper level.


## Setting up
There's a Makefile in the root directory that will build the Rush binary. You can run `make` to build it, it will create a `rush` binary in the root directory.

## Running
To run the Rush shell, you can run `./rush`.

## Commands
There are 2 types of commands in Rush: start commands and path commands.

### Standard commands
These commands are built into the shell, the goal of these commands is to provide augmented functionality to the shell. The full list of standard commands includes:
- `:e <VAR> <VAL>`: Prints the value of the environment variable `<VAR>` if `<VAL>` is not provided, otherwise sets the environment variable `<VAR>` to `<VAL>`.
- `:?`: Prints the about message.
- `:q`: Exits the shell.
- `:noop`: Does nothing.
- `:cd <DIR>`: Changes the current working directory to `<DIR>`.
- `:part <DISK_NAME>`: Prints the partitions on the disk `<DISK_NAME>`.
- `:omg`: I just had fun with the boys with this one.
- `:mem`: Allows to read memory of a process by it's PID.
- `:maps`: Prints the memory maps of the process.
- `:mountcron`: Mounts a special fuse filesystem (cronfs) that allows to read the contents of the crontab.
- `:unmountcron`: Unmounts the cronfs filesystem.

### Path commands
All other commands work the same as in other shells, the binaries are found in the `$PATH` and executed afterwards, passing the arguments to the binary.