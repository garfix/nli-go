# Quant functions

A quant is usually a relation 

    quant(quantifier, variable, range)
    
but there is also an ordered variant

    quant(quantifier, variable, range, order_function)
    
## Set order

Creates a new quant, based on an existing quant, but extended with an order function. If the original quant already had an order, it will be replaced.

    quant_order(Quant, OrderFunction, NewQuant)
    
* `Quant`: a `quant` relation
* `OrderFunction`: an atom, the name of an order relation
* `NewQuant`: a `quant` relation

If the quant is complex and contains sub-quants; then these will inherit the `OrderFunction` as well

    Example:
    
The order relation takes two entities and returns a negative number, 0, or a positive number. negative when E1 goes before E2, 0 when E1 has the same order position as E2, positive when E1 goes after E2.    
    
    by_easiness(E1, E2, R) :- if_then_else( cleartop(E1), unify(R, 1), unify(R, 0) );
    
    quant_order(Quant, by_easyness, NewQuant) 
