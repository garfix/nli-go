# Todo

- when listed values are used in the next sentence, they should be converted to multiple bindings

- create proper function diagrams for all processes
  - input
  - output
  - dependencies
  - side effects (which data stores are involved?)

- check agreement & create test

- The SHRDLU dialog does not have any definite anaphoric references (!) (where "the green block" refers to the block that was mentioned in the previous sentence) - so I will need an extra test in "relationships" that does just that

## Documentation

Describe documentation files per linguistic feature, with the following sections

- topic description (anaphora, multiple sentences, conjunctions, etc)
- examples that should be covered
- possible approaches
- nli-go's approach

This is how to write a book about it.

## Scripts and Frames

The restaurant script. Certain phrases invoke a script. (How?)

    "We went into a restaurant"

This phrase causes a number of discourse entities and relationships to be created at once, in the dialog context.

If the next sentence is "The waiter showed us our seat.", "the waiter" refers to an discourse entity in from the script. 

## The programming language "mentalese"

Make it consistent, complete, robust, etc. Have it conform existing paradigms.

- turn the "JSON" datatype into a "binary" datatype

- I must implement all entities with atoms. Currently they are variables, but it means that variables are used as values, and this is clumsy. Then there must be a mapping from these atoms to database ids.
- typed arguments
- operators > = [H|T]
- n-dimensional arrays as local variables
- extend a module with another module
* quant_foreach: add as second parameter the variable to which the ids must be bound

## Code

* binding set -> results / binding list
* relation set -> relation list
* better validation for built-in functions; especially multi-binding ones
* use functional programming: make all data immutable; use copy-on-write everywhere; stop making deep copies

## Long distance dependencies

The technique of ellipsis can possibly play a role in gapping / long distance relations too.

## Parsing source files

- check the correctness of all relations after parsing

## Interactive application

- interactive: arrow up/down for history

## The built-in mentalese application

- use parse tree as slot

## Performance

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


## Planning

- replanning: once a plan is being executed, it may need to be discarded and replanned, due to new circumstances
- stop / continue commands

## Quantifier Scoping

- Make "more than" "less than" work
- A range itself can contain quantified nouns (the oldest child in every family). The algorithm is not up to it. (See CLE)

## Stuff I'm not happy with

* the RelationTransformer; is only used in solutions, but should be removed from there as well, if possible
* `intent` moet een tag worden
* the entire function DialogContext::ReplaceVariable is bad; I should not change structs that are in the history list already  

## Interesting stuff

The catena concept may be interesting (for ellipsis a.o.):

    https://en.wikipedia.org/wiki/Catena_(linguistics)
