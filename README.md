
# go-quickswitch
Go Version of QuickSwitch (https://github.com/HellstromIT/quickswitch)

- [go-quickswitch](#go-quickswitch)
- [About](#about)
  - [WHY?](#why)
- [Installation](#installation)
  - [Download Packages](#download-packages)
  - [First run](#first-run)
  - [Add first search path](#add-first-search-path)
  - [Add shell functions](#add-shell-functions)
    - [Bash](#bash)
    - [Fish](#fish)
    - [zsh](#zsh)
    - [powershell](#powershell)

# About

## WHY?
The reasoning behind this package is the following. I've realised that I'm constantly switching between different git repositories on a normal day. To speed up the process of switching I wanted a tool that implemented would allow me to fuzzy search a set of directories containing git repositories. And that's how the idea behind quickswitch was born!:) 

# Installation

## Download Packages
Download the latest release from https://github.com/HellstromIT/go-quickswitch/releases, unpack it and install to your path:

Linux:
```
tar -xvf go-quickswitch_v0.1.0_linux_amd64.tar.gz
sudo cp go-quickswitch /usr/local/bin/
```

Mac:
(TBD)

Windows:
(TBD)

## First run

On first run the configuration file will be generated with the current working directory added as a search path. The program will exit telling you to re-run to search and also indicating where the config file is located.  


## Add first search path
There's no default search directories out of the box so to add a directory you wish to include run:

Linux/Mac:
```
go-quickswitch -add=/path/to/search/directory
```

Windows:
```
go-quickswitch.exe -add=/path/to/search/directory
```

The command will only search one level deep so if you have multiple levels that you wish to search you need to add them one at a time.


## Add shell functions
In order for the command to work correctly you will also need to create a function in your shell. The below functions assume that you want the command to be `qq`. 

After adding the relevant function restart your shell.

### Bash 
Add the following to your $HOME/.bashrc or in $HOME/.bash_functions.d/qq.sh

```
qq () {
  directory=$(qs)
  if [ -z "$directory" ]
  then
    echo
  else
    cd $directory
  fi
}
```

### Fish
Create a function in $HOME/.config/fish/functions/qq.fish

```
function qq
    set directory (qs)
    if set -q directory
        cd $directory
    else
        echo
    end
end
``` 

### zsh
Create a function in $HOME/.zshrc

```
qq () {
  directory=$(qs)
  if [ -z "$directory" ]
  then
    echo
  else
    cd $directory
  fi
}
```

### powershell
(TBD)

After adding the functions restart your shell and press qq<enter> to use the tool.

Happy Dir Switching!:)
