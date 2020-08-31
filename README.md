[![Go Report Card](https://goreportcard.com/badge/github.com/williamleven/BooGroCha)](https://goreportcard.com/report/github.com/williamleven/BooGroCha)
[![Build Status](https://travis-ci.com/sidusIO/BooGroCha.svg?branch=master)](https://travis-ci.com/siduIO/BooGroCha)
# BooGroCha
A CLI client for chalmers' group room booking system.

## Features
- [x] Book group rooms
- [x] Integration Chalmers booking system
- [x] Save personal settings in config file
- [x] List available bookings for selected time
- [x] Delete bookings
- [x] List bookings
- [x] Relative times for bookings
- [x] [AUR package](https://aur.archlinux.org/packages/boogrocha)
- [x] Rooms sorted base on user preferences
- [ ] [Option to only show rooms from preferred campus](https://github.com/williamleven/BooGroCha/issues/6)
- [ ] [Integrate with Chalmers Library booking system](https://github.com/williamleven/BooGroCha/issues/7)
- [x] Book rooms in the Johanneberg Student Union building
- [x] Prompt for password if not set


## Installation

### Compile from source

#### Requirements
* go (>= 1.12)
* $GOPATH/bin in $PATH

#### Steps
```bash
$ git clone https://github.com/williamleven/BooGroCha
$ cd BooGroCha/cmd/bgc
$ go install
```

### [AUR package](https://aur.archlinux.org/packages/boogrocha)

## Usage

### Book a room
Allows for booking rooms by showing the available rooms for the given date and time.

```bash
$ bgc book <date> <time>
```
* **\<date\>** can be either an absolute date (`YYYYMMDD`, `YYMMDD`) or a relative date (`MMDD`, `DD`, `D`, `today`, `tomorrow`, `monday`, `tuesday` ...)
* **\<time\>** can be either (`HH:mm-HH:mm`, `HH-HH`, `H-H`, `H-HH`, `HH-H`) or aliases like `lunch`

The `book` sub-command also takes the following optional flags:
* `--cid <cid>` to perform the book request with a specified cid
* `--campus <campus>` or `-c <campus>` to filter the available rooms by a specified campus. (Valid campuses are `J`, `Johanneberg`, `L` and `Lindholmen`. Note that these values are not case-sensitive)
* `--size <size>` or `-s <size>` to filter the available rooms by size and will only show the rooms that are big enough. (When a size is specified the list of available rooms will also show the capacity of each room)

### List booked rooms

```bash
$ bgc list
```

### Delete booked rooms

```bash
$ bgc delete
```

### Configuration
Allows the user to set parameters in a config file.

#### Setting a variable
```bash
$ bgc config set <variable> <value>
```
* **\<variable\>** can for example be `cid` or `pass`
* **\<value\>** should be the value that you want to set the variable to (NOTE: when setting the password you will be prompted for input instead of setting it directly)


#### Showing a variable
```bash
$ bgc config get <variable>
```
* **\<variable\>** can for example be `cid` or `campus` (NOTE: You cannot show the `pass` variable this way)
