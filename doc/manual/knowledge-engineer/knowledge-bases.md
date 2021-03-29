# Knowledge bases

To solve a query (in the form of a relation set), the system uses several knowledge bases. These knowledge bases are specified in the config file.

Each knowledge base provides bindings for the active variables of the query. Different knowledge bases are used at different moments in the problem solving process.

## Fact base

The fact base can be simply a file of relations, a MySQL database, or a SPARQL datastore.

It has a number of "mappings" that map the input to fact base specific output. It maps one or more relations from the input to one or more relations from the output. The bindings come from the fact base.

## Function base

A function base processes relations that normally take the form of a function.

For example

    join(Result, Sep, String1, String2, ...)

## Rule base

A rule base contains "Prolog" rules that expand a consequent to a sequence of antecedents.

    sister(A, B) :- female(A) female(B) parent(A, C) parent(B, C)

A rule base treats an consequent relation as a goal, and executes its antecedents.

## Aggregate base

An aggregate base provides functions that take a set of bindings as input and perform function on the all the bindings of

An example of an aggregate function is

    count(X, N)

In the function Bind, two things happen:

* the number of different value of X in bindings is calculated
* all bindings are extended with a value for N

## Nested structure base

A nested structure base passes control to the child relations of a relation.

Currently these are the nested structures:

    quant()
    sequence()
    not()
    call()

