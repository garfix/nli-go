# List functions

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
