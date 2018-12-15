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

You can use the executable as you would use any command-line application. It has two sub-commands:

Use nli to answer a question, based on a configuration stored in a JSON config file. It returns a JSON string with the answer and / or an error.

```
./nli answer fox/config.json "Did the quick brown jump over the lazy dog?"
```

Or use it to suggest the next words the user can type.

```
./nli suggest fox/config.json "Did the quick"
```

The config file is described [here](doc/manual/config.md).
