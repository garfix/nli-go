# Todo

- go:list_make()
- go:list_deduplicate
- go:list_sort
- go:list_index
- go:list_length
- go:list_get

- local variables (mutable)
- functions calls for arguments
- typed arguments
- operators
- n-dimensional arrays as local variables

- check of alle relaties goed zijn bij het inlezen van de source files
- extend a module with another module
- interactive: arrow up/down for history

* quant_foreach: add as second parameter the variable to which the ids must be bound 
* agreement, especially for number, because it reduces ambiguity (reintroducing feature unification?)
* syntax check while parsing: is the number of arguments correct?
* remove the square brackets where they are not needed
* SparqlFactBase: todo predicates does not contain database relations (just ontology relations), so this needs to be solve some other way
* clarification questions must be translatable (they must go through the generator)
* use relations as functions (with special role for the last parameter as the return value)
* write a good tutorial
* think of a better replacement to make_and() to an "and" sequence 

* (?) to_list(E1)
    collect all distinct values of E1 into a list, replace the value of E1 in all bindings with this list; remove duplicate bindings
    not yet needed; maybe postpone
    
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

## Aggregation functions on bindings

`number_of`, `exists`, and the functions that still need to be build, `min` and `max`, work on bindings, and it is better to make this explicit.

    bindings_max(E)
    
for example. On the other hand, I could make a single function

    to_list(E)
    
that converts the binding variables into a list. And create

    list_max(E)
    
But I don't like `to_list` because it must change the variable E and this is against the rules in Prolog. If it wouldn't change E then the number of bindings would stay unnecessary large.            

## Misc

* Separate interfaces (api) from implementations (model)
* Blocks World examples

## Rules

Test if this works or make it work. Create a stack of current relations to be solved, and check if the stack already contains the bound relation.

    married_to(A, B) :- married_to(B, A);
    
* Allow the dynamically added rules to be saved (in the session).
* Specify which predicates a rule base allows to be added.    

## Syntax

- Perhaps replace the syntax of functions like number_of(N, X) to
    number_of(X: N)
    join('', firstName, lastName: name)
    join('', firstName, lastName -> name)
    name = join('', firstName, lastName)
- should you be allowed to mix predicates of several sets? Is this confusing or a necessity to keep things manageable?
- Must be able to write whword in place of whword(); but wait, maybe we need multiple variables as well?

## Aggregations

- Add min, max

## Relations

Find a way to ensure completeness of information about all relations used in a system. An interpretation should not even be attempted if not all conversions have a chance to succeed.

* convert number words into numbers

# Multiple languages

- Introduce a second language

# Quantifier Scoping

- Make "more than" "less than" work
- A range itself can contain quantified nouns (the oldest child in every family). The algorithm is not up to it. (See CLE)
