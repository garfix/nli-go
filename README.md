# NLI-GO

nli-go is small executable program, written in Go, that provides a Natural Language Interface to databases. It is a semantic parser and execution engine. I use it to experiment with nli techniques. It is not stable yet and backward-incompatible changes will be made from time to time.

## Demo

A demo of this library can be found [here](http://patrickvanbergen.com/dbpedia/app/). It allows you to use a few sentences to query DBPedia. 

## Purpose

This library helps a developer to create a system that allow end-users to use plain English / French / German to interface with a MySQL or SPARQL database. This allows your the system to handle interactions like this:

    > Q How many children had Lord Byron? 
    > A: He had 2 children.
    
    > Q: Was Michael Jackson married to Elvis Presley's daughter?
    > A: Yes

Every part of the system is configurable.

## Techniques

Some of the techniques used:

* An Earley parser to create a syntax tree from an input sentence, with semantic attachments
* Mentalese, a based internal language, based on Predicate Logic, to process the user input
* A Prolog-like language for rule based reasoning
* Support for Sparql (DBPedia) and MySQL as well as an in-memory data stores
* Using data from multiple databases in a single request
* A dialog context to remember information from earlier in the conversation
* Anaphora resolution: the use of pronouns and other references to earlier entities
* Generalized quantifiers
* The distinction between classic negation (`not`) and strong negation (`-`)
* A generator to produce human readable responses
* Modules and namespaces, for modular development

## Build the nli executable

NLI-GO is a command-line application called "nli". It's written in Go and you must compile it into an executable for your OS.

You can download and install GO from [here](https://golang.org/dl/)

From the root of NLI-GO build the executable with

    go build -o bin/nli bin/nli.go
    
The executable is now available as `bin/nli`. You can add the extension .exe and move the executable to another location if you like.    

## Run the executable with the sample applications

NLI-GO comes with some sample applications, located in the directory "resources". In this example you tell "Hello World" to the hello world application:

    bin/nli -c resources/helloworld "Hello World"    

and it responds with

    Welcome!

This is the response of the application, or the error, if something went wrong. If you need more control over the output of the system, you can add `-r json`; like this

    bin/nli -c resources/helloworld -r json "Hello World"    
  
and NLI-GO responds with a JSON string like this:

    {
        "Success": true,
        "ErrorLines": [],
        "Productions": [
            "Anaphora queue: [] ",
            "TokenExpression: [Hello World] ",
            "Parse trees found: : 1 ",
            "Parser: [s [hello Hello] [world World]] ",
            "Relationizer: go_intent(start_conversation) dom_hello() ",
     ...
            "Answer Words: [Welcome!] ",
            "Answer: Welcome! "
        ],
        "Answer": "Welcome!",
        "OptionKeys": [],
        "OptionValues": []
    }
    
If the system responds with a clarification question, it does this with a number of options the user can choose from

* OptionKeys: the keys of these options
* OptionValues: the values of these options

And if you want to specify a session identifier to allow NLI-GO to resolve back-references to earlier parts of the dialog, use `-s` with an identifier of your choice.     

    bin/nli -c resources/helloworld -s 64152 "Hello World"    
    
## Trying it out

If you want to experiment with NLI-GO, copy one of the application directories in `resources` and make changes to it.

## Docs

Much information on how to build an NLI-GO application can be found in [How to build an NLI application](doc/manual/knowledge-engineer/index.md).

If you want to follow my thoughts as I develop NLI-GO, you can read it here: [My personal log](doc/remarks.md)

And this is an overview of [How NLI-GO processes a request](doc/manual/system-developer/processing.md).
