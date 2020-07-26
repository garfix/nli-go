# Todo

todo

* list as term 
* if_then_else(if, then, else)
    if `if` succeeds, then `then`, else `else`
* quant_order(quant, sort_func, ordered_quant)    
* quant_ordered_list(quant, List, orderFunc)
    goes through compound quant, using orderFunc, placing results in List 
* list_order(List1, sort_func, List2)
    sorts List1 by sort_func, places result in List2
* list_foreach(List, E, scope)
    instantiates List in E, and executes scope
---
* do -> quant_foreach(quant, E, scope) // note! single quant per predicate, use min      
* find -> quant_check(quant, E, scope) // note! single quant per predicate


* agreement, especially for number, because it reduces ambiguity
* syntactic placeholders for `sem(n)` (`$np1`)
* namespaces for relations: `:find()`, `db:support()`
* should boolean functions have P1 as argument? different or for read/write?
* check if the nested functions are called correctly
* syntax check while parsing: is the number of arguments correct?
* do not allow zero valued predicates in the grammar
* SparqlFactBase: todo predicates does not contain database relations (just ontology relations), so this needs to be solve some other way
* entity type (multiple) inheritance
* sortal restrictions (using predicates.json and adding 'parent' to entities.json)
* agreement checking (reintroducing feature unification)
* clarification questions must be translatable (they must go through the generator)
* use relations as functions (with special role for the last parameter as the return value)

* (?) to_list(E1)
    collect all distinct values of E1 into a list, replace the value of E1 in all bindings with this list; remove duplicate bindings
    not yet needed; maybe postpone

## Agreement

    'boy', sense: block(E), agr(E, number, 1)
    'boys', sense: block(E), agr(E, number, multiple)
    
    'pick' 'up', sense: pick_up(E1) agr(E1, number, 1) // first person singular
    'pick' 'up', sense: pick_up(E1) agr(E1, number, multiple) // plural
    
    'pick' 'up', sense: pick_up(E1) number(E1, multiple)
    > number's second argument must have a single value; declare this in some way

## Syntactic sugar

    sense: find([sem(1) sem(4)], marry(P1, E1, E2))
    sense: find(a b, marry(P1, E1, E2))
    sense: find(#1 #2, marry(P1, E1, E2))
    sense: find($1 $2, marry(P1, E1, E2))
    { rule: np_check(P1) -> np(E1)$1 marry(P1) 'to' np(E2)$2,                    sense: find_all([$1 $2], marry(P1, E1, E2)) }
    { rule: np_check(P1) -> np(E1)$subject marry(P1) 'to' np(E2)$object,         sense: find_all([$subject $object], marry(P1, E1, E2)) }
    { rule: np_check(P1) -> np(E1) marry(P1) 'to' np(E2),                        sense: find_all([$np1 $np2], marry(P1, E1, E2)) }

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

Test if this works or make it work:

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
- Constants like "all", are they universal, or english?

# Quantifier Scoping

- Make "more than" "less than" work
- A range itself can contain quantified nouns (the oldest child in every family). The algorithm is not up to it. (See CLE)
