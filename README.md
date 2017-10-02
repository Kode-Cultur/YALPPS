# YALPPS #

is **yet another LAN party port scanner**. The idea behind this
project is to speed up the setup phase of LAN parties by providing a
simple-to-use binary that will check whether any of the TCP or UDP
ports required by a game are blocked by firewall rules and informing the
user of the findings.

## Architecture ##

A `YALPPS` server is run from one computer, together with a list of
TCP and UDP ports that need to accessible for the games that will be
played. From each other computer in the network, a client is started.
The client connects to that server and receives instructions on when
to test whether which port is open. After that, the user is presented
with a list of ports (for both inbound and outbound connections) that
are currently blocked by firewall rules.

## Compiling ##

This project uses the Go language, currently in version 1.9.
For information on how to install and setup Go, [refer to the official
documentation](https://golang.org/doc/install).

Make sure to place this repository somewhere in your `$GOPATH`.

To fetch the required dependencies, run:

    go get -t .

Then, to compile the project, simply run:

    go build

## Licensing ##

YALPPS is licensed under the Apache 2.0 license. See `LICENSE` for the
complete license text.
