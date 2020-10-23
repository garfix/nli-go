# List functions

## List make

Creates a new list, based on existing variables with multiple bindings. The values of `X`, `Y` and `Z` are added to the list, and they are removed from the bindings. The resulting bindings are deduplicated.

    go:list_make(List)
    go:list_make(List, X, Y, Z, ...)
    
* `List`: a variable that will contain the list
* `X`, `Y`, `Z`: variables with one or more values

Example: 

When calling `go:list_make(List, X)` with bindings 

    [{X: 2, Y: 1}{X: 3}{}] 

the resulting bindings will be

    [{List: [2, 3, 1]}]
    
As you can see the values are added, first in order of argument appearance, and second in order of binding appearance. First `X` is added as `2` and `3`; then `Y` is added as `1`.

## List order

Creates a new list, based on an existing list, but ordered by an order function.

    list_order(List, &OrderFunction, NewList)
    
* `List`: a list
* `OrderFunction`: a reference to a rule that functions as an order function
* `NewList`: a variable (that will be bound to the ordered list)

The order relation takes two entities and returns a negative number, 0, or a positive number. negative when E1 goes before E2, 0 when E1 has the same order position as E2, positive when E1 goes after E2.    
    
    by_easiness(E1, E2, R) :- if_then_else( cleartop(E1), unify(R, 1), unify(R, 0) );
    
    list_order(List, &by_easiness, NewList)

## List foreach

Binds each of the values of list to `Variable`, and executes `Scope` for each value.

    list_foreach(List, Variable, Scope)
    
* `List`: a list
* `Variable`: a variable
* `Scope`: a relation set

## List deduplicate

Creates a new list, based on existing list, but with all duplicate elements removed. The order of the values will not change.

    go:list_deduplicate(List, NewList)
    
* `List`: a list
* `NewList`: a variable that will contain a list
