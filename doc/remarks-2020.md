# 2018-01-18

I thought about anaphora resolution (handling pronouns like "it", and phrases like "the block", that refer to entities
that were just previously mentioned in the conversation).

This is how I intent to deal with it:

I will create a queue of entity-references that represent recently mentioned things in the conversation. This queue
consists of key-cabinet entries. Each entry holds the id's of the entity in one or more knowledge bases.

Whenever an entity is "consciously" processed by the system, it will be added to the queue. The queue will have a
certain maximum size and the last references will be removed when the queue is full.

In my solution I will not try to "detect" anaphoric references in noun phrases. I will treat anaphoric references just
like normal noun phrases. But I will make a change.

The critical structure is found in the processing of quantifications (quants). This process first fetches the id's of
all entities in the range. Then it processes the scope. And finally it filters the result through the quantifier.

The addition I make is in the loading of the range. When the range is specified by the word "he", for example, its sense
consists just of "male(E1)". This means that the entities considered would be all male persons. I will not load all
entities and filter them with the keys from the anaphoric queue. In stead, I will first attempt to restrict the range
with the available keys in the queue.

An example:

Anaphoric queue

    [{db=db2, id = 9591, type = block}]
    [{db=db1, id = 312, type = person} {db=db2, id = 111, type = person}]
    [{db=db1, id = 8, type = person} {db=db2, id = 9012, type = person}]
    [{db=db1, id = 31001, type = block}]

Input: When did the red man marry?

    when() 
    quant( 
        the()               // quantifier 
        red(E1) man(E1)     // range
        marry(E1, E2))      // scope

When the quant is processed, the processor will take the range

    "man(E)" 

and take the first entry of the anaphoric queue

    [{db=db2, id = 9591, type = block}]

since this doesn't match (type = block, not person), it tries the next

    [{db=db1, id = 312, type = person} {db=db2, id = 111, type = person}]

This gives a match.

Only if no match was found using the queue, will the range be matched with the knowledge bases without a key constraint.

This is the basic idea; I expect there will need to made some adjustments to make this work.

# 2020-01-11

I put a new version online http://patrickvanbergen.com/dbpedia/app/ that allows case insensitive names. This will reduce
the number of failed queries quite a lot.

Did you know?

    Who is <name>?  --> Asks for a description 
    Who is <description>? --> Asks for a name 

# 2020-01-05

I made case-insensitive names possible, and at the same time checking the database in the parsing process. Introduced
s-selection to help finding the names. s-selection restricts predicate arguments and this in turn narrows the search
space for proper nouns in the database.

# 2020-01-02

Happy new year! 

I am introducing semantics in the parsing process, because I need some semantics to determine the type of the entity in
a name.

I want to use the relationizer that I already have for this, but it is too much linked to the nodes that I generate
after the parse is done.

Now I just had an interesting idea: what if I do the sense building as part of the chart states. That way, when the
parse is done, I just need to filter out the widest complete states and I will have the complete sense ready, without
having to create a tree and then parse that tree.
