# Knowledge engineer

A knowledge engineer does not develop an application but writes the grammar, the rules and the solutions.

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

Now run the test. NLI-GO will tell you what words can't be parsed:

    ERROR: Incomplete. Could not parse word: What
    
The form of the word is the "surface form" of the word itself, in lowercase.

The part-of-speech is the syntactic category of a word. If you don't know the part-of-speech run it in the online parser of CoreNLP:

    http://corenlp.run/
    
It will give you the syntactic categories of the words, as codes. The description of these codes can be found here:

    https://www.ling.upenn.edu/courses/Fall_2003/ling001/penn_treebank_pos.html
    
You can use CoreNLP's notation or the one I use mostly, as given in syntax.md. Any form will do, as long as it is compatible with the grammar.

The sense is the "meaning" of the word, in a relational form. For most words it is simply this:

    sense: isa(E, find)
    
This example belongs to the verb "find". The sense merely states that any find object (entity, or for a verb: event) belongs to the class `find`.

This is the basic sense. It is simple to find because you just copy the word form into the sense.
There are a number of exceptions to this basic rule, however.

Any derivative of a verb must have the infinitive as its sense:

    finding     =>  find(E)
    found       =>  find(E)
    
If the application needs it you can specify voice and aspect:

    finding     =>  find(E) aspect(E, progressive) 
    found       =>  find(E) tense(E, past)
    
But if the application does not need it, you can leave them out.

## Extend the grammar

You may get the error

    ERROR: Incomplete. Could not parse word: What

This is because the existing grammar rules are not enough to parse this sentence. You must add the rules that allows NLI-GO to parse the sentence. 

To determine what rewrite rule to add, you need to follow the workings of the parser. Start with s() and perform a number of rewrites until the part-of-speech of the new word can be matched. 

You are free to write grammar rules. If you can make use of existing grammars, that is preferable to writing your own. 

But most grammars are incomplete so you will likely need to be creative at times.  

I have used the rules from the following book mostly

    The structure of modern English - Brinton
    
A grammar also has has relations, just like the lexical entry. These "senses" form the domain-independent representation of the meaning of the sentence.
This means that the meaning stays close to the syntax and far away from domain and database representations.

As a simple example, here is a typical declarative sentence:

    rule: s(P1) -> np(S1) vp(P1),       sense: subject(P1, S1)
    
As you can see, P1 is used both on the left side and on the right side. This variable represents an event, and this event is passed on to the child vp().

The np() forms the subject of the sentence, and this is a meaningful part so its sense is declared. 
The variable can be chosen at will, but I use P for predication, S for subject and O for object. 

## Extend the solutions

At some point the system will say

    ERROR: There are no solutions for this problem

Then you need to add a solution to the solutions file. Here is an example that contains all sections of a solution:

~~~
    {
        condition: question(_) how_many(B) have_child(A, B),
        transformations: [
            have_child(A, B) => have_n_children(A, Number);
        ],
        responses: [
            {
                condition: exists(),
                preparation: gender(A, Gender),
                answer: gender(A, Gender) have_child(A, C) count(C, Number)
            }
            {
                answer: dont_know()
            }
        ]
    }
~~~

The condition specifies what the input relations must look like for this solution to apply.

The transformations transform some relations from the input to the set that will actually be processed. This is kind of
"interpretation" on the part of the system. (When the user says A he actually means B). Transformations are optional.

A solution can have several responses, that depend on certain conditions. So each response has an optional condition.

Sometimes you need some extra information in the answer that was not retrieved in the question itself. For example you
may want to to answer "He had 2 children", whereas the gender of the subject was not part of the question. "preparation"
allows you to fetch this extra information.

"answer" are just a passive set of relations passed to the generator. They are not processed in any way. All variable bindings must already be available at this point.

A special function of "answer" is "make_and().

    answer: root(R) name(A, N) make_and(A, R)

This construction allows you to create a nested structure of AND's, so that you can respond with

    John, Kale and Louis
