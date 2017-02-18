Welcome to Golang Examples
===================

I publish my small and functional projects. I hope it will help.
> **Note:**

> - All examples developed on localhost.
> - İf you want to work on real world, you have to implement the server-side codes to server which have **static IP address**.


--Chat and Comminication (clients beetwen server) 
-------------

First, this example explains the relationship between a server and clients which are using **Websocket**.
Simply, **Websocket** provides full-duplex communication.
For more information about [websocket][1]


###Server-side source code files
#### <i class="icon-file"></i> main.go 
It starts to run server and listen the port for handling messages and connections.

#### <i class="icon-file"></i> packages/Client.go 
It explains how clients which connected a server behave on the side of server.

#### <i class="icon-file"></i> packages/Server.go 
It explains how server behaves to itself and clients.

#### <i class="icon-file"></i> packages/Types.go 
Some variables and types for comminication.

  [1]: https://en.wikipedia.org/wiki/WebSocket

