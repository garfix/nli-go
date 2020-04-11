# Knowledge engineer

The knowledge engineer writes the grammar, the rules, the solutions, and the mapping to the database.

## Before you start

Make sure you have a working resource directory with the basic configuration files (grammar, solution, etc) that you can
extend.

## The content writing process

There may be different approaches of how to teach the program how to process a number of sentence. The method describes here is iterative. 

* Collect a sentence from an end user
* Add an automated test 
* Extend the grammar
* Extend the solutions
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

## Extend the grammar

Now run the test. NLI-GO will throw an error, like for example this one:

    ERROR: Incomplete. Could not parse word: What
    
This means you will have to add rewrite rules to the grammar. Check the other grammars for examples on how to write a
grammar. Do not copy complete grammars. Copy just the single lines you need. This way you will be able to comprehend
your grammar.

More on the grammar you can find [here](entity-grammar.md)

## Extend the solutions

At some point the system will say

    ERROR: There are no solutions for this problem

This means that the system does not not how to handle the input. For each type of sentence there is a separate solution.

More on solutions you can find [here](solution.md)
