![a peer to peer TCP tic tac toe game in plain Golang](https://github.com/Abdenasser/tcp-tac-toe/blob/main/imagetcp-tac-toe.jpg?raw=true)


# TCP Tic Tac Toe

A peer to peer TCP tic tac toe game in plain Golang.

## How to play

make sure you have go installed on your machine.

1. clone the repo `git clone https://github.com/Abdenasser/tcp-tac-toe.git`

2. cd into the repo `cd tcp-tac-toe`

3. run the server `go run server.go`

4. once the server is running, open two terminal windows and connect to the server using `telnet localhost 8080`

5. start playing!

## playing over the internet

if you want to play with a friend over the internet, you can use a service like [ngrok](https://ngrok.com/)
to expose your local server to the internet.

1. run the server `go run server.go`

2. run ngrok `ngrok tcp 8080`

3. from nother machines, connect to the ngrok address using `telnet <ngrok_address> <ngrok_port>`

4. start playing!