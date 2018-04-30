# Zuper-UDP
[![Build Status](https://travis-ci.com/mohamedmahmoud97/Zuper-UDP.svg?token=aQtpZzy2UuNChYAfpRmS&branch=master)](https://travis-ci.com/mohamedmahmoud97/Zuper-UDP)
[![Go Report Card](https://goreportcard.com/badge/github.com/mohamedmahmoud97/Zuper-UDP)](https://goreportcard.com/report/github.com/mohamedmahmoud97/Zuper-UDP)
[![Documentation Status](https://readthedocs.org/projects/zuper-udp/badge/?version=latest)](http://zuper-udp.readthedocs.io/en/latest/?badge=latest)
[![API Reference](
https://camo.githubusercontent.com/915b7be44ada53c290eb157634330494ebe3e30a/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f676f6c616e672f6764646f3f7374617475732e737667
)](https://godoc.org/github.com/mohamedmahmoud97/Zuper-UDP)
[![License](https://img.shields.io/dub/l/vibe-d.svg)](https://github.com/mohamedmahmoud97/Zuper-UDP/blob/v2/LICENSE)

 A library for reliable data transfer service on top of the UDP protocol that supports loadbalancing and a lite version of reverse proxy. It also simulates packet loss and packet corruption and handles these issues for reliability.
 
 <p align="center">
  <img width="600" height="299" src="https://github.com/mohamedmahmoud97/Zuper-UDP/blob/master/images/loadbalancer.jpg">
</p>

## Installation
You can use `go get` to install to your `$GOPATH`, assuming that you have a github.com parent folder already created under `src`:
```
$ go get -t github.com/mohamedmahmoud97/Zuper-UDP
```

## Usage
```
$ go run examples/loadbalancer/loadbalancer.go
 
 _ _ _                          _    _ ___   _ __
|_ _  | _   _ _ __   ___  _ __ | |  | |._ \ | ._ \
   /  /| | | | ._ \ / _ \| '__|| |  | || \ || |_) |
  /  /_| |_| | |_) | |_)/| |   | |__| ||_/ || .__/
 /_ _ _|_ _ _| .__/ \_ _ |_|   |_ __ _|___/ |_|    v1.0  LOADBALANCER
             |_|

 A server-side loadbalancer with udp reliable data transfer protocol

 #####################################################################

```

## Test it 

You can test the library with running any of the bash files: `gbn.sh` for testing GoBackN algorithm, `sr.sh` for SelectiveRepeat, and `sw.sh` for Stop-and-Wait. You can also change the parameters for the client, the server, and the loadbalancer in `device_info/client.in` and `device_info/server.in`, and `device_info/loadbalancer.in`. If you want to neglect the simulation probability packet loss just change last attribute in `server.in` to 0.

You can also test it by running a loadbalancer by running `go run examples/loadbalancer/loadbalancer.go` and run some servers and clients and make the server address attribute for all clients is the address of the loadbalancer.

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
127.0.0.1            The IP address of server
10066                Well-known port number for server.
5                    Maximum sending sliding-window size (in datagram units).
5                    Random generator seedvalue.
0.3                  Probability p of datagram loss (real number in the range [ 0.0 , 1.0 ]).
```

#### **For Loadbalancer**
```
127.0.0.1:10000                The IP address of loadbalancer
127.0.0.1:10001                The IP address of first server
127.0.0.1:10101                The IP address of second server
127.0.0.1:10201                The IP address of third server
127.0.0.1:10301                The IP address of fourth server
127.0.0.1:10401                The IP address of fifth server
```

#### **test run of server**
[![asciicast](https://asciinema.org/a/GSjxa39RyY1SgRVkHj7wCntSf.png)](https://asciinema.org/a/GSjxa39RyY1SgRVkHj7wCntSf)

#### **test run of client**
[![asciicast](https://asciinema.org/a/U9RGjOWl1sKWf3qkxEVMyE5AI.png)](https://asciinema.org/a/U9RGjOWl1sKWf3qkxEVMyE5AI)

## Contributing
Feel free to contribute to this project. Either for adding features, found a bug, .. etc.
