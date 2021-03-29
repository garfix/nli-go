# Todo

- functions calls for arguments
- typed arguments
- operators
- n-dimensional arrays as local variables

- check of alle relaties goed zijn bij het inlezen van de source files
- extend a module with another module
- interactive: arrow up/down for history

## Performance
- if the system is instantiated just to process messages, dont't install all language components; lazy load; specially for all grammars
- on the other hand: maybe the rules take longer to parse, and they are always necessary
- create a service? (a stay resident application that processes messages)

- let => var
- document: these are two ways of doing a child stack frame
- the internal factbase is inefficient; for every new and removed fact, all facts are matched
- would be nice to have `is_first()` and `is_last()`: a check if the current binding is the first / last of the active bindings; such a function takes both a single binding and all bindings as input

## Blocks
The animation also reveals another problem: when the system builds a stack, it first decides on a location, then builds it. When building the first block, it may need to place the objects on top of it in some location. And it chooses the exact location where the stack should be. Later, the rest of the stack is still placed there. A solution could be to exclude this intended location from free space.

* database mappings: allow a rule to be used only for given sorts; for performance
* binding set -> results / binding list
* relation set -> relation list  
* better validation for built-in functions; especially multi-binding ones
* quant_foreach: add as second parameter the variable to which the ids must be bound 
* agreement, especially for number, because it reduces ambiguity (reintroducing feature unification?)
* syntax check while parsing: is the number of arguments correct?
* SparqlFactBase: todo predicates does not contain database relations (just ontology relations), so this needs to be solve some other way
* clarification questions must be translatable (they must go through the generator)
* use relations as functions (with special role for the last parameter as the return value)
* write a good tutorial
* think of a better replacement to make_and() to an "and" sequence 
* change rewrite rules from categories with variables to relations (see also Generator)

## generation of multiple entities

Replace `make_and()` by a `make_list()` and add list unification syntax

    { rule: entities(E1) -> entity(A) ',' entities(Tail),                         condition: go:unify(E1, [A _ _ | Tail]) }
    { rule: entities(E1) -> entity(A) 'and' entities(B),                          condition: go:unify(E1, [A B]) }
    { rule: entities(E1) -> entity(A),                                            condition: go:unify(E1, [A]) }
    { rule: entities(E1) -> entity(E1) }    

The last one is used with just a single constant.

## Agreement

    'boy', sense: block(E), agr(E, number, 1)
    'boys', sense: block(E), agr(E, number, multiple)
    
    'pick' 'up', sense: pick_up(E1) agr(E1, number, 1) // first person singular
    'pick' 'up', sense: pick_up(E1) agr(E1, number, multiple) // plural
    
    'pick' 'up', sense: pick_up(E1) number(E1, multiple)
    > number's second argument must have a single value; declare this in some way 

## Stuff I'm not happy with

* the RelationTransformer; is only used in solutions, but should be removed from there as well, if possible

## Rules

Test if this works or make it work. Create a stack of current relations to be solved, and check if the stack already contains the bound relation.

    married_to(A, B) :- married_to(B, A);
    
* Allow the dynamically added rules to be saved (in the session).
* Specify which predicates a rule base allows to be added.    

## Syntax

- Perhaps replace the syntax of functions like number_of(N, X) to
    count(X: N)
    join('', firstName, lastName: name)
    join('', firstName, lastName -> name)
    name = join('', firstName, lastName)
- should you be allowed to mix predicates of several sets? Is this confusing or a necessity to keep things manageable?
- Must be able to write whword in place of whword(); but wait, maybe we need multiple variables as well?

## Relations

Find a way to ensure completeness of information about all relations used in a system. An interpretation should not even be attempted if not all conversions have a chance to succeed.

* convert number words into numbers

# Multiple languages

- Introduce a second language

# Quantifier Scoping

- Make "more than" "less than" work
- A range itself can contain quantified nouns (the oldest child in every family). The algorithm is not up to it. (See CLE)
