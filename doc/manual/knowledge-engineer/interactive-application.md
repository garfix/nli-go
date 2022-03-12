# Building an interactive application

The "blocks" web demo shows how the system can be used as a client-server application, where the client is a Javascript web-app and the server is formed by calls to nli-go.

The main idea is that the client and the server communicate by sending messages to each other. The client sends a message, the server processes it and sends a message back. The client responds to it until the message is empty.

When the user types a sentence, the client sends this message (in JSON) to the server:

    [
        {
            "predicate":"go_tell",
            "arguments":[
                {
                    "type":"string",
                    "value":"Find a block which is taller than the one you are holding and put it into the box."
                }
            ]
        }
    ]

The message simply consists of a list of relations and in this case just the single relation `go:tell()`. The relations are processed by nli-go's function `System:SendMessage`. The relation / procedure `go:tell` is located in `respond.rule` and looks like this:

    go:tell(Input) :-
        go:create_goal(
            go:respond(Input)
        );

`go:tell` doesn't really do anything. It just creates a goal: respond to input, and adds it to the list of goals of the dialog context.

Next, `System:SendMessage` continues to excecute the active goals. Each goal is executed by a process. When no process exists, one is created. In this case a process is created that consists of the relation set `go:respond(Input)`.

This process just runs and runs, until it meets a relation `go:wait_for()`. This relation contains a check. If the check succeeds, `wait_for` just continues; but when it doesn't, the process stops (pauses), and the server. The most `wait_for` looks like this:

    go:wait_for(
        go:print(Uuid, :Output)
    )

It's a funny check, it's more like a command: "print!". But it really checks if the relation `go:print(Uuid, :Output)` exists. This check is send to the client as follows:

    [
        {
            "predicate": "go_print",
            "arguments": [
                {
                    "type": "string",
                    "value": "6FF779CC6BA7EE91"
                },
                {
                    "type": "string",
                    "value": "OK"
                }
            ]
        }
    ]

The client picks up the check, interprets it as a command, and executes it, by printing the text specified (here: "OK"). It then sends back the message: I have done what you asked, you may assume the check is true:

    [
        {
            "predicate":"go_assert",
            "arguments":[
                {
                    "type":"relation-set",
                    "set":[
                        {
                            "predicate":"go_print",
                            "arguments":[
                                {
                                    "type":"string",
                                    "value":"6FF779CC6BA7EE91"
                                },
                                {
                                    "type":"string",
                                    "value":"OK"
                                }
                            ]
                        }
                    ]
                }
            ]
        }
    ]

As you can see, this message again is just a relation set that the server just executes:

    go:assert(
        go:print("6FF779CC6BA7EE91, "OK")
    )

The server, `System:SendMessage`, executes this relation and then continues to execute the process. This time when the `wait_for` comes along, the check is true, and the process continues.



