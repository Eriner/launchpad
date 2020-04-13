# launchpad [![GoDoc](https://godoc.org/github.com/Eriner/launchpad?status.svg)](https://godoc.org/github.com/Eriner/launchpad)
A package allows you to talk to your Novation Launchpad X in Go. Light buttons or read your touches.

Provides a state machine and middleware!

This library was originally a fork of [rakyll/launchpad](https://github.com/rakyll/launchpad) but has been completely rewritten. This library only supports the Launchpad X and provides additional features, including a grid state machine.

~~~ sh
go get github.com/eriner/launchpad
~~~

Portmidi is required to use this package.

```
$ apt-get install libportmidi-dev
# or
$ brew install portmidi
```

## Usage
An example has been heavily commented in the cmd/main.go file.
