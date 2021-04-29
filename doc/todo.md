# Todo

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
- If you hold block A and are told to put block A in the box (or on something), don't put it down first (don't clear hand)

## Documentation

* write a good tutorial

## Code

* binding set -> results / binding list
* relation set -> relation list
* better validation for built-in functions; especially multi-binding ones

## relations that I no longer use

- go:isa(E, Sort)

## The programming language "mentalese"

Make it consistent, complete, robust, etc. Have it conform existing paradigms.

- mutable variables now have global scope; this is really wrong and should be fixed => scope must be limited to declaring rule
- if_then, if_then_else => if
- let => var
- functions calls for arguments
- typed arguments
- operators > = [H|T]
- keywords if/then  
- n-dimensional arrays as local variables
- extend a module with another module
* quant_foreach: add as second parameter the variable to which the ids must be bound
* use relations as functions (with special role for the last parameter as the return value)

## Rules

Test if this works or make it work. Create a stack of current relations to be solved, and check if the stack already contains the bound relation.

    married_to(A, B) :- married_to(B, A);
    
* Allow the dynamically added rules to be saved (in the session).
* Specify which predicates a rule base allows to be added.    
* change rewrite rules from categories with variables to relations (see also Generator)

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

# Quantifier Scoping

- Make "more than" "less than" work
- A range itself can contain quantified nouns (the oldest child in every family). The algorithm is not up to it. (See CLE)

## Stuff I'm not happy with

* the RelationTransformer; is only used in solutions, but should be removed from there as well, if possible
