## Vim for go

### Installing on MacOS

The prefer way to install is with `brew` loke so:

    $ brew install vim

Plugin https://github.com/fatih/vim-go

After installing run `:GoInstallBinaries`

Or use https://vim-bootstrap.com/ to generate a start point for `.vimrc`

## Coc for go

Add Coc plugin

```
" Use release branch (recommend)
Plug 'neoclide/coc.nvim', {'branch': 'release'}
```

Install extensions
```
:CocInstall coc-json coc-tsserver coc-go coc-python
```

https://github.com/josa42/coc-go

## Entr. Event Notify Test Runner

A utility for running arbitrary commands when files change. Uses kqueue(2) or inotify(7) to avoid polling. entr was written to make rapid feedback and automated testing natural and completely ordinary.

ref: [https://github.com/eradman/entr](https://github.com/eradman/entr)

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

### Tips

* [For go dev](https://dev.to/jogendra/using-vim-for-go-development-5hc6)
* [Vim cheatsheet](https://vim.rtorr.com/)

### iTerm

https://gist.github.com/nobitagit/729fc16b8c16edb9a2fe390d6f312c66
