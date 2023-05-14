Each system has up to two processes: a language process and an action process. The language process processes a request, may send some clarification requests, may start up an action process, and returns without waiting for the action process to complete. The language process goes so far as to plan the action, but it does not itself execute it.

The action process receives a plan and executes it.

A system can have only one of these processes active at a time.

As long as a language request is active, each user request will be interpreted as a response to a clarification request. If no language request is active, the request will create a new language process. A request like "stop" is applied to the currently running action process.

A request that tries to start an action process, while an action process is already active, will fail. No two actions can be done at the same time. A response could be "I can't do that. I have to complete my current action first"

When the last running procedure has finished, the system will send an "all clear" message. An automated procedure will take this message as a trigger to start the next request.
