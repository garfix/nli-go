## Quant foreach

Find all entities specified by a `quant` (minimally), assign each of them in turn to a variable and execute `Scope`.

Fails as soon as a scope returns no results. 

    go:quant_foreach(Quant ..., Scope)
    
* `Quant`: a quant
* `Scope`: a relation set    

Check [quantification](quantification.md) for more information.

## Quant check

Find all entities specified by `Quants`, check if the number of entities that pass `Scope` is the same as specified by the quantifier of `Quant`. 

    go:quant_check(Quants, Scope)
    
* `Quants`: one or more quants
* `Scope`: a relation set      

Check [quantification](quantification.md) for more information.

## Quant to list

Creates a new quant, based on an existing quant, but extended with an order function. If the original quant already had an order, it will be replaced.

    go:quant_ordered_list(Quant, &OrderFunction, List)
    
* `Quant`: a `quant` relation
* `OrderFunction`: a reference to a rule that functions as an order function
* `List`: a variable (to contain a list)

If the quant is complex and contains sub-quants; then these will be ordered by the `OrderFunction` as well

    Example:
    
The order relation takes two entities and returns a negative number, 0, or a positive number. negative when E1 goes before E2, 0 when E1 has the same order position as E2, positive when E1 goes after E2.    
    
    by_easiness(E1, E2, R) :- if_then_else( cleartop(E1), unify(R, 1), unify(R, 0) );
    
    go:quant_ordered_list(Quant, &by_easyness, List) 
