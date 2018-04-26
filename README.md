# Zuper-UDP
[![Build Status](https://travis-ci.com/mohamedmahmoud97/Zuper-UDP.svg?token=aQtpZzy2UuNChYAfpRmS&branch=master)](https://travis-ci.com/mohamedmahmoud97/Zuper-UDP)
[![Go Report Card](https://goreportcard.com/badge/github.com/mohamedmahmoud97/Zuper-UDP)](https://goreportcard.com/report/github.com/mohamedmahmoud97/Zuper-UDP)
[![API Reference](
https://camo.githubusercontent.com/915b7be44ada53c290eb157634330494ebe3e30a/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f676f6c616e672f6764646f3f7374617475732e737667
)](https://godoc.org/github.com/mohamedmahmoud97/Zuper-UDP)
[![License](https://img.shields.io/dub/l/vibe-d.svg)](https://github.com/mohamedmahmoud97/Zuper-UDP/blob/v2/LICENSE)

 A library for reliable data transfer service on top of the UDP protocol.

## Installation
You can use `go get` to install to your `$GOPATH`, assuming that you have a github.com parent folder already created under `src`:
```
$ go get -t github.com/mohamedmahmoud97/Zuper-UDP
```

## Usage
```
$ go run examples/server/server.go sw
 
 _ _ _                          _    _ ___   _ __
|_ _  | _   _ _ __   ___  _ __ | |  | |._ \ | ._ \
   /  /| | | | ._ \ / _ \| '__|| |  | || \ || |_) |
  /  /_| |_| | |_) | |_)/| |   | |__| ||_/ || .__/
 /_ _ _|_ _ _| .__/ \_ _ |_|   |_ __ _|___/ |_|    v1.0  SERVER
             |_|

 A server-side udp reliable data transfer protocol

 #####################################################################

started connection in server on port :10066 ... 

```

## Test it 

You can test the library with running any of the bash files: `gbn.sh` for testing GoBackN algorithm, `sr.sh` for SelectiveRepeat, and `sw.sh` for Stop-and-Wait. You can also change the parameters for the client and the server in `device_info/client.in` and `device_info/server.in`

#### **For Client**
```
127.0.0.1            IP address of server
10066                Well-known port number of server.
127.0.0.1            IP address of client.
10044                Port number of client.
test.pdf             Filename to be transferred.
5                    Initial receiving sliding-window size (in datagram units).
```

#### **For Server**
```
10066                Well-known port number for server.
5                    Maximum sending sliding-window size (in datagram units).
5                    Random generator seedvalue.
0.3                  Probability p of datagram loss (real number in the range [ 0.0 , 1.0 ]).
```

#### **test run of server**
[![asciicast](https://asciinema.org/a/GSjxa39RyY1SgRVkHj7wCntSf.png)](https://asciinema.org/a/GSjxa39RyY1SgRVkHj7wCntSf)

#### **test run of client**
[![asciicast](https://asciinema.org/a/U9RGjOWl1sKWf3qkxEVMyE5AI.png)](https://asciinema.org/a/U9RGjOWl1sKWf3qkxEVMyE5AI)

## Contributing
You can feel free to contribute to this project. Either for adding features, found a bug, etc.


