## LSP

Enable Language Server Protocol

### Format on save

Lua configuration

    vim.cmd [[autocmd BufWritePre <buffer> lua vim.lsp.buf.formatting_sync()]]
    -- or
    vim.cmd [[autocmd BufWritePre * lua vim.lsp.buf.formatting_sync()]]

### Install LSP Server

Each server has to be installed manually or with
[lsp-installer](https://github.com/williamboman/nvim-lsp-installer)

_gopls_

    $ go install golang.org/x/tools/gopls@latest
