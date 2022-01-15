# Todo

## The programming language "mentalese"

Make it consistent, complete, robust, etc. Have it conform existing paradigms.

- I must implement all entities with atoms. Currently they are variables, but it means that variables are used as values, and this is clumsy. Then there must be a mapping from these atoms to database ids.
- functions calls for arguments
- typed arguments
- operators > = [H|T]
- n-dimensional arrays as local variables
- extend a module with another module
* quant_foreach: add as second parameter the variable to which the ids must be bound
* use relations as functions (with special role for the last parameter as the return value)

## Solutions

Call solutions "intents", just like everybody else.

## Code

* binding set -> results / binding list
* relation set -> relation list
* better validation for built-in functions; especially multi-binding ones
* use functional programming: make all data immutable; use copy-on-write everywhere; stop making deep copies

## Anaphora

- Donkey sentences: "If Pedro owns some donkey, he beats it." "Some donkey" creates a discourse entity; and "it" refers to it. (Also: "Pedro owns a donkey. He beats it.") 
- indefinite descriptions ("Jones owns a Porsche. It fascinates him.")
- Store the senses of the entities that go into the anaphora queue, for later matching
- Forward references
- The referent of an anaphoric expression sometimes can be found only by using world knowledge ("John beat Peter. He started to cry.")
- Some antecedents depend on the syntactic role ("John supports Peter. He admires him.")
- Winograd schema challenge https://en.wikipedia.org/wiki/Winograd_Schema_Challenge

## Long distance dependencies

The technique of ellipsis can possibly play a role in gapping / long distance relations too.

## Parsing source files

- check the correctness of all relations after parsing

## Interactive application

- interactive: arrow up/down for history

## The built-in mentalese application

- use parse tree as slot

## Performance

- if the system is instantiated just to process messages, dont't install all language components; lazy load; specially for all grammars
- on the other hand: maybe the rules take longer to parse, and they are always necessary
- create a service? (a stay resident application that processes messages)

- the internal factbase is inefficient; for every new and removed fact, all facts are matched
- would be nice to have `is_first()` and `is_last()`: a check if the current binding is the first / last of the active bindings; such a function takes both a single binding and all bindings as input

## Database

* database mappings: allow a rule to be used only for given sorts; for performance
* SparqlFactBase: todo predicates does not contain database relations (just ontology relations), so this needs to be solve some other way

## Blocks demo

The animation also reveals another problem: when the system builds a stack, it first decides on a location, then builds it. When building the first block, it may need to place the objects on top of it in some location. And it chooses the exact location where the stack should be. Later, the rest of the stack is still placed there. A solution could be to exclude this intended location from free space.

- When the demo is done. Do it in German as well, as proof of multilinguality.
- Add an autoplay function for demo mode; show characters in response one-by-one  
- If you hold block A and are told to put block A in the box (or on something), don't put it down first (don't clear hand)

## Documentation

* write a good tutorial

## relations that I no longer use

- go:isa(E, Sort)

## Rules

Test if this works or make it work. Create a stack of current relations to be solved, and check if the stack already contains the bound relation.

    married_to(A, B) :- married_to(B, A);
    
* change rewrite rules from categories with variables to relations (see also Generator)

## Relations

Find a way to ensure completeness of information about all relations used in a system. An interpretation should not even be attempted if not all conversions have a chance to succeed.

* convert number words into numbers

## Language features

- Parse multiple sentences in a single line of input. 

Agreement (see 2021-12-28)

An example of the syntax I will use for feature structures and unification:

    { rule: vp(P1, E1) -> np(E1) aux_be(_) tv_gerund(P1, E1, E2) np(E2),            agree: number(P1, E1) }
    { rule: noun(E1) -> 'blocks',                                                   tag: number(E1,  plural) }


## Quantifier Scoping

- Make "more than" "less than" work
- A range itself can contain quantified nouns (the oldest child in every family). The algorithm is not up to it. (See CLE)

## Stuff I'm not happy with

* the RelationTransformer; is only used in solutions, but should be removed from there as well, if possible
* the async routines are too complex

## Interesting stuff

The catena concept may be interesting (for ellipsis a.o.):

    https://en.wikipedia.org/wiki/Catena_(linguistics)
