# Knowledge engineer

The knowledge engineer writes the grammar, the rules, the solutions, and the mapping to the database.

## Before you start

Make sure you have a working resource directory with the basic configuration files (grammar, solution, etc) that you can
extend. Copy the directory `helloworld` for instance. Look at [this page](config.md) for the contents of the `config.json` file.

## Relations

In this system you convert sentences in natural language to a semantic representation in the form of relations. Relations, "senses", or "predications" are extremely versatile. To get a feeling for this concept, check [this](relation.md).

## The content writing process

There may be different approaches of how to teach the program how to process a number of sentence. The method describes here is iterative. 

* Collect a sentence from an end user
* Add an automated test 
* Extend the grammar
* Extend the solutions
* Create the database mappings
* Extend the generations
* Repeat until done

## Collect a sentence

You need to ask the future user of the system what kind of questions he/she would ask. 

Once you have a running system that allows the user to enter one or more questions, you can also log their questions in order to find out what they expect from the system.

Do not impose sentences on the user and expect them to use these; every new user will be frustrated that the natural form does not work.

## Add an automated test

Create an automated test for all the new question and its expected answer.

This allows you to go through the following phases quickly without having to enter the question again and again.

When you try to add a new question, it is possible that you break earlier question / answers. So you need to run the previous tests as well.

## Change the regular expression for the tokenizer

The standard tokenizer expression creates tokens of all adjoining wordlike-characters: "apple", "11", "c64". All other characters each get a separate token: "?", "'", "-". This works fine for simple words, but if you need email addresses or other complex forms as single tokens, you need to create a custom regular expression. 

You can name the regular expression used to tokize a sentence if you are not satisfied with the standard expression. For example

    "tokenizer": "([_0-9a-zA-Z]+|[^\\s])"   

## Extend the grammar

Now run the test. NLI-GO will throw an error, like for example this one:

    ERROR: Incomplete. Could not parse word: What
    
This means you will have to add rewrite rules to the grammar. Check the other grammars for examples on how to write a
grammar. Do not copy complete grammars. Copy just the single lines you need. This way you will be able to comprehend
your grammar.

Each grammar rule turns a phrase structure into a single word, or some other phrase structures. It may also create a semantic attachment, called "sense". The system will combine these senses to create the meaning of the sentence.  

More on the grammar you can find [here](entity-grammar.md)

Creating a grammar by hand is not easy. The best way to create one is to keep it as simple as you need for the sensences you need to support, and only abstract when necessary. 

In the past few years I found out how to turn different types of phrases into semantic structures. These are described in [Creating a grammar](creating-a-grammar.md).   

## Extend the solutions

At some point the system will say

    ERROR: There are no solutions for this problem

This means that the system does not not how to handle the input. For each type of sentence there is a separate solution.

More on solutions you can find [here](solution.md)

## Create a database mapping

Some relations are resolved by the database. A relation like `parent(X, :18)` will be transformed by the system to a query like `SELECT parent_id FROM person WHERE child_id = 18`. Each relation is converted to a simple query. The system does not create complex queries.   

To use this you must

* Configure the database and the tables you need in `config.json`
* Create a `.map` file

A database can be just a `.relation` file, or it can be a MySQL or Sparql database. In most cases you will just want to read from the database. If you also want to write to the database, you will need to create a separate `write` map-file.

See [knowledge bases](knowledge-bases.md) for the different types of databases and other knowledge bases available.

The map file defines how a relation in your application is mapped to one or more relations in the database. 

For Sparql databases you will also need to specify how a relation's predicate is mapped to a uri. See names.json for example in dbpedia/db.

## Common sense reasoning

The goals that have been set by the input sentence and the solution, may also be reached by following inference rules. These rules allow you to create procedures that solve the problem in a common sense way. See [Common sense reasoning](common-sense-reasoning.md) for more on this. 

## Generate a response

There is a separate grammar ("write grammar") for responses, because answers are very different in structure from questions. And where read-grammars produce semantic attachments ("sense"), write grammars take existing semantics as conditions for the sentence to be generated.

More about response generation [here](generation.md).
