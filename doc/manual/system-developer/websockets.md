The system has a client-server architecture. A request by a user may take a longer time to execute and all this time the server is active. A request may need multiple processes, and these are implemented as goroutines.

The interaction of the client with the server is implemented as a websocket. The clients sends messages to the server. The server sends messages to the client.

These are messages sent by the server:

- print (print a sentence)
- choose (have the user make a decision)
- processlist_clear (no running processes: feel free to interact in any way)

These are messages sent by the client:

- respond (send user input)
- acknowledge (in response to print)
- chosen (in response to choose)

There are also system-specific messages.
