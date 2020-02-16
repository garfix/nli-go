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

Use nli to answer a question, based on a configuration stored in a JSON config file. It returns a JSON string with the answer and / or an error.

```
./nli -c "../resources/blocks/config.json" "Pick up the box"
```

the response could be:

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
* Productions: debug lines to help you debug
* Answer: the actual answer to wanted

If the system responds with a clarification question, it does this with a number of options the user can choose from

* OptionKeys: the keys of these options
* OptionValues: the values of these options

The config file is described [here](doc/manual/knowledge-engineer/config.md).
