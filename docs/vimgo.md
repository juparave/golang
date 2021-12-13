## Vim for go

Plugin https://github.com/fatih/vim-go

After installing run `:GoInstallBinaries`


## Entr. Event Notify Test Runner

A utility for running arbitrary commands when files change. Uses kqueue(2) or inotify(7) to avoid polling. entr was written to make rapid feedback and automated testing natural and completely ordinary.

ref: https://github.com/eradman/entr

Command to work with go:

    $ <<< 'dumbooctopus.go' entr -cc go run main.go dumbooctopus.go

* `dumbooctopus.go`: is the file being watched for changes
* `entr`: the command
* `-cc`: to clear the screen
* `go run main.go dumbooctopus.go`: the program to run when the file changes


### Splits

For a split window: You can use `Ctrl-w +` and `Ctrl-w -` to resize the height of the current window by a single row. For a vsplit window: You can use `Ctrl-w >` and `Ctrl-w <` to resize the width of the current window by a single column. Additionally, these key combinations accept a count prefix so that you can change the window size in larger steps. [e.g. `Ctrl-w 10 +` increases the window size by 10 lines]

* To resize all windows to equal dimensions based on their splits, you can use `Ctrl-w =`.
* To increase a window to its maximum height, use `Ctrl-w _`.
* To increase a window to its maximum width, use `Ctrl-w |`.


