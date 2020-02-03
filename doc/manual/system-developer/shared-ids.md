# Shared ids

Shared ids are an important concept. It requires some introduction.

The system starts by turning a natural language question sentence into a relation set.
Then this relation set is matched against a series of knowledge bases in order to find the answer to the question.
The question, in its essence, consists of a series of variable bindings.
One these bindings are found, the system ends by creating an answer sentence for the user.

A user sentence may contain proper names, such as "Lord Byron" or "Ada Lovelace".
In a simple system such a name is an identifier for a single person in a database.
When information about the entities is distributed over multiple databases, this situation changes.

## Multiple databases

The system described here may contact multiple databases to solve a single user question.

For example: let's take the user question:

How old was Albert Einstein when he graduated?

given that database 1 contains the information:

person(18, "Albert Einstein")
birth(18, "1958-08-29")

and database 2 contains the information

name(r22, "Albert Einstein")
graduation(r22, "1896-10-03")

The notable point here is that Albert Einstein has different identifiers in both databases.

Yes, this example is a bit contrived, but it nails the problem on the head. We have a single question whose answer requires information from two databases, that both identify the person at hand in a different way.

## Multiple identifiers

If some variable in the relation set (say E3) contains the identifier of the person, the relation set will never match.
It cannot match both 18 (database 1) and r22 (database 2).

_Shared ids_ are a solution to this problem. A shared id is an id that does not occur in any database, yet database ids
are mapped to it via mapping tables (db-id <--> shared id)

In this example:

    database-1:
    db-id = 18, shared-id = aeinstein
    
    database-2:
    db-id = r22, shared-id = aeinstein

The shared ids are not part of the database. Thet are stored in JSON files and added to the fact base structure. The
name of the JSON file is stored in the `sharedIds` attribute of a fact base.

All types of fact bases have this property, so it can even be used to join
disjoint data that is stored in distinct databases such as a MySQL database and a Sparql database.

## How does the system use shared ids?

The system uses shared-ids everywhere in the system. As soon as a proper noun is located in a database, its shared id is
looked up and bound to a local variable. The database id is left in the database.

Only when the system accesses a database (fact base) directly it will swap shared ids to database local ids just before
making the call to the database. But as soon as the call returns the local ids will be swapped back to shared ids.

When a shared id is expected but not found, the system will show an exception.

## When to use shared ids

Only when you have entity types that are located in multiple databases with distinct ids, you will need shared ids.
If this is not the case, the system will simple use the local database ids as shared ids.

## How do I create the shared-id JSON files?

By hand. Seriously. It can be a lot of work. It may be possible to use heuristics like "an e-mail address always
uniquely identifies a person", or "most objects in each database have an EAN". But there will possibly be exceptions
which need to be handled by hand. The positive side is that for the most part, this is a one-time job.
