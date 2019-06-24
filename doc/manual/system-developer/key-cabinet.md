# Key cabinet

The key cabinet is a difficult but important concept. It requires some introduction.

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

The key cabinet is a solution to this problem. It stores bindings of a variable, per database.

In this example:

E3 { database-1: 18, database-2: r22 }

The key cabinet is passed all through the solution process, together with the bindings. But while the bindings change all the time, the key cabinet stays the same.

And only when a particular database is accessed, is a value retrieved from the key cabinet and used to query the database.
It is never bound to any variable and does not enter the bindings.

