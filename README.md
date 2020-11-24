
# go-quickswitch
Go Version of QuickSwitch (https://github.com/HellstromIT/quickswitch)

- [go-quickswitch](#go-quickswitch)
- [About](#about)
- [Installation](#installation)
  - [Linux:](#linux)
  - [Windows:](#windows)
  - [Mac:](#mac)
  - [Bash](#bash)
  - [Fish](#fish)
  - [zsh](#zsh)
  - [powershell](#powershell)



# About

# Installation
Download the latest release from https://github.com/HellstromIT/go-quickswitch/releases, unpack it and install to your path:

## Linux:
```
tar -xvf go-quickswitch_v0.1.0_linux_amd64.tar.gz
sudo cp go-quickswitch /usr/local/bin/
```

Create configuration file (this is a bit convoluted at the moment. Will be fixed in the future):

```
echo '{"Directories": []}' > ~/.config/quickswitch.json
```

## Windows:
(TBD)

## Mac:
(TBD)

There's no default search directories out of the box so to add a directory you wish to include run:

```
go-quickswitch -add=/path/to/search/directory
```


The command will only search one level deep so if you have multiple levels that you wish to search you need to add them one at a time.

In order for the command to work correctly you will also need to create a function in your shell. The below functions assume that you want the command to be `qq`. 

After adding the relevant function restart your shell.

## Bash 
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

## Fish
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

## zsh
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

## powershell
(TBD)

After adding the functions restart your shell and press qq<enter> to use the tool.

Happy Dir Switching!:)
