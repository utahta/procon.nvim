# procon.nvim

This is an experimental plugin for Neovim to automate the typical tasks in competitive programming.

## Supported websites are
- [AtCoder](https://atcoder.jp/)

## Supported languages are
- C++14/gcc
- C++17/gcc

## Installing

You'll need to install [Go](https://golang.org/dl/).

##### vim-plug
```viml
Plug 'utahta/procon.nvim', {'do': 'make'}
```

## How to use

1. Enter the `:ProconSet` command to set up the site, contest, etc. that you challenge.
  Then the plugin will create a new file with an appropriate name in the current directory.
  If there is already a file with the same name, you can enter any name you like.
1. Write the code.
1. Enter the `:ProconCheck` command to make sure your program outputs the correct answer.
  If your program outputs the correct answer, you'll get the AC. And of course, if it outputs the wrong answer, you'll get the WA.
  Note: If the task accepts multiple answers, you may get the WA even if your program outputs the correct answer, because the plugin only compares it to the sample cases.
1. If you are confident that you'll get the AC, enter the `:ProconSubmit` to submit your code.
