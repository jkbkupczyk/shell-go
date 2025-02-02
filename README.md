# shell-go

A minimal POSIX compliant shell, implemented in Go, that's capable of interpreting 
shell commands, running external programs and builtin commands like:
echo, type, pwd, cd and more. Includes stdout and stderr redirection and autocompletion.

**Note**: This project is a part of ["Build Your Own Shell"](https://app.codecrafters.io/courses/shell/overview) from [codecrafters.io](https://codecrafters.io).

## Features

- REPL
- running builtin commands [exit, echo, type, pwd, cd]
- running external program with arguments
- support for non quoted, single quoted, double quoted arguments
- stdout and stder redirection
- stdout and stderr appending
- autocompletion
- signal handling (`Ctrl+C`, `Ctl+D`)

## Installation

Before you do any of the steps below, make sure you have [Golang](https://go.dev) installed.

1.  Clone the repo
    ```sh
    git clone https://github.com/jkbkupczyk/shell-go.git
    ```
2.  Build and run
    ```sh
    go build -o bin/myshell ./cmd/myshell/.
    ./bin/myshell
    ```
3.  Run tests (optional!)
    ```sh
    go test -v -timeout 30s ./...
    ```

You can also use [Make](https://www.gnu.org/software/make) commands: `run` - for building and running application and `test` for running tests - see [Makefile](Makefile).


## Roadmap / TODOs

- more tests
- refactors, fixes
- piping
- history
- variable interpolation
- job control
- new functionalities (reference [Bash Reference Manual](https://www.gnu.org/software/bash/manual/bash.html))

## Contributing

Feel free to create issues or submit pull requests!
