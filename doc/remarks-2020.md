# 2020-01-29

I will interrupt my quest to answer question 5 to make a release that handles linked databases. I will make a new
integration test that handles the following question:

    Which poem wrote the grandfather of Charles Darwin?

Charles Darwin and the link to his grandfather are stored in db1. Erasmus Darwin and his poem is stored in db2. Both
databases are not directly linked, but there will be shared id tables.

# 2020-01-27

Shared-id == meta-id.

I will first un-use the key cabinet, and when that works, I will give it the new function of providing the shared-id -
db-id mapping.

I must make a test for the multiple-database query.

# 2020-01-25

What if a mapping between databases would exist? What would it look like? Preferably there would be a shared identity.

Of course many entities already have public id. Books have ISBN numbers. People have email addresses and
Social security numbers. Maybe they can be used in a certain domain, maybe not. But that's up to the designer of the
system. It is possible in some cases.

If I have 4 databases, A B C D, I could map A to B, A to C, A to D, B to C, B to D, C to D, but I could also map A to
shared and from shared to B.

    db 1: 8119              =>  johnbrown-1
    db 2: 23                =>  johnbrown-1

When I find an id in database A I can then do two things:

A) I could fetch the id's in A B etc. and place them in the key cabinet.
B) I could fetch the shared id and assign it to the query variable

In both cases I would need to find the database specific ID just before querying.

I prefer B. I would not need a key cabinet any more.

In most cases of course there is only one database. In this case the shared id is identical to the database id.

When and how should you create the mapping? Can it be done on-the-fly or must it be done periodically?

The mapping can be made as follows: suppose db1 has person fields `[id, email, name]` and db2 has fields `[id2, last
name, first name, email]` then the mapping should be created off-line by going through the persons in db1 and matching
them to the persons in db2, via heuristics or hard identities. The result would be mapping table for each database.

Even if not all db's have the entity, there must still be a shared id.

What would change?

- the key cabinet can go
- the id-term is a shared id; which can default to a database id, when there is only a single database
- in the configuration I would need to know which entity types have shared ids, and for each database a mapping table
- when a database is queried, the mapping from shared id to db id must be made; the response must be mapped again to the
    shared id

# 2020-01-24

Do I want a meta-id? (an entity that links one or more actual database id's)

Do I want to extend the notion of the id-term to meta-id? That if you create an id-term, that you are really creating a
meta-id with an initial single db id?

Can this replace the key cabinet?

meta-id: {db1: 15, db2: 18}

Matching meta's:

    {} + {db1: 16} => {db1: 16}
    {db1: 15} + {db1: 16} => no 
    {db1: 15} + {db2: 16} => {db1: 15, db2: 16}
    {db1: 15, db2: 4} + {db2: 4 db3: 108} => {db1: 15, db2: 4, db3: 108}

I can still use the key cabinet. The id term will then be a key from the cabinet. The cabinet holds the database ids;
possibly along with the entity type.

But how do these meta's behave?

When the sentence holds the name of an entity, the system can ask the user to clarify which of the entities, possibly in
different databases, is meant. A meta id will have one or more db ids.

But what happens further in the processing of the sentence? As long as a variable is not matched with a database id, it
will not be bound to a meta id, or perhaps only to an empty meta id. Once there's a database match, the db id will be
added to the meta id.

After that, the meta id will not gather any more ids. It is not possible that {db1: 13} will be used in `parent({db1:
13}, E2)` for db2. Because we do know which db2 entity matches 13 in db1, and we cannot leave it out entirely, because
the first argument _is known_, only not for db2, and we don't want any values from db2 as if the argument hadn't been
bound yet, for this would yield disallowed values. It would also not be possible to find new values for E2 in db2. So
this jeopardizes the idea that the system is capable of linking separate databases.

A way out would be a mapping function that links the entities of separate databases, but there is no such thing.

# 2020-01-23

In order to resolve anaphoric references I need to store the id's of earlier entities. The point is that the id's at the
moment are not unique even in a database, in the case of number ID's. An ID may belong to different tables. I can do two
things:

- store the table name in the id
- keep track of entity types during the processing.

Wait, I can look up the entity type for each relation, so this is not a problem.

---

I just noticed that my solver binds variables to ids of specific databases. Actually it is meant to bind to some
database-independent ids, that have links to ids in separate databases. This will be a problem when entities will be
found in multiple databases, as I intend to do. This is what I built the key cabinet for.

# 2020-01-22

I made the dialog context much simpler and more extendable by removing the fact base.

# 2020-01-18

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

This is the basic idea; I expect there will need to be made some adjustments to make this work.

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
