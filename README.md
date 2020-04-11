# launchpad [![GoDoc](https://godoc.org/github.com/Eriner/launchpad?status.svg)](https://godoc.org/github.com/Eriner/launchpad)
A package allows you to talk to your Novation Launchpad X in Go. Light buttons or read your touches.

Provides a state machine and middleware!

This library is a fork of [rakyll/launchpad](https://github.com/rakyll/launchpad) and has legacy support for the [Launchpad S (Green-Red Launchpads)](https://www.amazon.com/dp/B00E3XTKAG).

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
An example has been heavily commented in the cmd/helo_x/main.go file.

### Coordinate system

The coordinate system is illustrated below.
~~~
+--------- arrow keys -----------+  +--- mode keys ---+
{0, 8} {1, 8} {2, 8} {3, 8} {4, 8} {5, 8} {6, 8} {7, 8} | ableton
----------------------------------------------------------------
{0, 0} {1, 0} {2, 0} {3, 0} {4, 0} {5, 0} {6, 0} {7, 0} | {8, 0} vol
----------------------------------------------------------------
{0, 1} {1, 1} {2, 1} {3, 1} {4, 1} {5, 1} {6, 1} {7, 1} | {8, 1} pan
----------------------------------------------------------------
{0, 2} {1, 2} {2, 2} {3, 2} {4, 2} {5, 2} {6, 2} {7, 2} | {8, 2} sndA
----------------------------------------------------------------
{0, 3} {1, 3} {2, 3} {3, 3} {4, 3} {5, 3} {6, 3} {7, 3} | {8, 3} sndB
----------------------------------------------------------------
{0, 4} {1, 4} {2, 4} {3, 4} {4, 4} {5, 4} {6, 4} {7, 4} | {8, 4} stop
----------------------------------------------------------------
{0, 5} {1, 5} {2, 5} {3, 5} {4, 5} {5, 5} {6, 5} {7, 5} | {8, 5} trk on
----------------------------------------------------------------
{0, 6} {1, 6} {2, 6} {3, 6} {4, 6} {5, 6} {6, 6} {7, 6} | {8, 6} solo
----------------------------------------------------------------
{0, 7} {1, 7} {2, 7} {3, 7} {4, 7} {5, 7} {6, 7} {7, 7} | {8, 7} arm
----------------------------------------------------------------
~~~

## Demo: Light your touchs

A simple program to light every touch:

## License
    Copyright 2020 Matt Hamilton. All Rights Reserved.
    Copyright 2013 Google Inc. All Rights Reserved.
    
    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at
    
         http://www.apache.org/licenses/LICENSE-2.0
    
    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
