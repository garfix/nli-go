# NLI-GO

nli-go is a library, written in Go, that provides a Natural Language Interface to databases. I use it to experiment with
nli techniques. It is not stable yet and backward-incompatible changes can be made at any time.

## Demo

A demo of this library can be found [here](http://patrickvanbergen.com/dbpedia/app/). It allows you to use a few sentences to query DBPedia. 

## Purpose

This library helps a developer to create a system that allow end-users to use plain English / French / German to interface with a database. That means that an end user can type a question like

    > Q How many children had Lord Byron? 
    > A: He had 2 children.
    
    > Q: Was Michael Jackson married to Elvis Presley's daughter?
    > A: Yes

Every part of the system is configurable.

## Techniques

Some of the techniques used:

* An Earley parser to create a syntax tree from an input sentence, with semantic attachments
* Mentalese, a based internal language, based on Predicate Logic, to process the user input
* A Datalog (a very basic Prolog) implementation for rule based reasoning
* Support for Sparql (DBPedia) and MySQL as well as an in-memory data stores
* Linking data from multiple databases in a single request
* A dialog context to remember information from earlier in the conversation
* Anaphora resolution: the use of pronouns and other references to earlier entities
* Generalized quantifiers
* A generator to produce human readable responses

## Docs

Documentation is located in the docs directory, here you can find:

* [My personal log](doc/remarks.md)
* [The processing of a request](doc/manual/system-developer/processing.md)

