# Solution route

A solution route is one way to solve a problem. It consists of relation groups, together with a kb id and a cost of execution.

A solution route is always ordered by least expensive first.

## The cost of a solution

Solving a problem using database lookup takes time. The order of processing the relations is important. The most restricting relations must be processed first.

Some relations are present in more than one knowledge base. Also, they may need to be grouped together in order to be transformable into target knowledge base relations.

Solution routes are a means to do this. Relations are grouped in RelationGroups and a solution route is an array of these groups.

A relation group contains relations that can be transformed together, a database id and a cost of execution.

Solution routes; different ways of solving a single problem:

    [(relations, db-id, cost), (relations, db-id, cost), (relations, db-id, cost)]
    [(relations, db-id, cost), (relations, db-id, cost), (relations, db-id, cost)]
