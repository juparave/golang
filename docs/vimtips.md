# Vim Tips

[Vim Cheatsheet. from devhints](https://devhints.io/vim)

- `F3` open nerdtree
- `gg=G` reindent whole file

## Commenting Multiple Lines

Follow the steps given below for commenting multiple using the terminal.

- First, press `ESC`
- Go to the line from which you want to start commenting. Then, press `ctrl v`, this will enable the visual block mode. Use the down arrow to select multiple lines that you want to comment.
- Now, press `shift I` to enable insert mode.
- Press `#` or `//` and it will add a comment to the first line. Then press `ESC` and wait for a second, `#` or `//` will be added to all the lines.

## Uncommenting Multiple Lines

- Press `ctrl v` to enable visual block mode.
- Move down and select the lines till you want to uncomment.
- Press `x` and it will uncomment all the selected lines at once.

## Resize splits

For a split window:

`Ctrl-w +` and `Ctrl-w -` to resize the height of the current window

For a vsplit window:

`Ctrl-w >` and `Ctrl-w <` to resize the width of the current window

To resize in larger steps

`Ctrl-w 10 +`

You can use the `:resize` command or its shortcut `:res` to change the height of the window. To change the height to 60 rows, use:

    :resize 60

You can use `:vertical` resize to change the width of the current window. To change the width to 80 columns, use:

    :vertical resize 80
