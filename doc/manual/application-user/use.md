# Use

## Install the program

Build this application's executable:

```
cd ~/go/src/nli-go/app/cli
go build nli.go
```

If you like, you can move it to a place where it can be found from any location. In a Linux environment you might use:

```
sudo mv nli /usr/local/bin
```

## Command-line use

You can use the executable as you would use any command-line application.

The following call uses `nli` to answer the question "What does the box contain?", based on a configuration stored in a
JSON config file and storing its dialog context (spanning multiple sentences) in a session file in `./sessions/123.json`

```
./nli -s 123 -c "../resources/blocks/config.json" "What does the box contain?"
```

the response is a JSON string, that looks like this:

~~~
{
    "Success":true,
    "ErrorLines":[],
    "Productions":[
        "Anaphora queue: []",
        "Tokenizer: [Pick up the box]",
        ...
        "Answer: OK"
    ],
    "Answer":"OK",
    "OptionKeys":[],
    "OptionValues":[]
}
~~~

* Success: has the sentence been processed completely?
* ErrorLines: in case of an error, tells you what went wrong
* Productions: progress information to help you debug
* Answer: the actual answer to wanted

If the system responds with a clarification question, it does this with a number of options the user can choose from

* OptionKeys: the keys of these options
* OptionValues: the values of these options

The application writes its log files in the directory `./logs`

The config file is described [here](../knowledge-engineer/config.md).
