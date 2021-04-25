# List predicates

## Make list

Creates a new list, based on existing variables with multiple bindings. The values of `X`, `Y` and `Z` are added to the list, and they are removed from the bindings. The resulting bindings are deduplicated.

    go:make_list(List)
    go:make_list(List, X, Y, Z, ...)
    
* `List`: a variable that will contain the list
* `X`, `Y`, `Z`: variables with one or more values

Example: 

When calling `go:make_list(List, X)` with bindings 

    [{X: 2, Y: 1}{X: 3}{}] 

the resulting bindings will be

    [{List: [2, 3, 1]}]
    
As you can see the values are added, first in order of argument appearance, and second in order of binding appearance. First `X` is added as `2` and `3`; then `Y` is added as `1`.

## List append

Creates a new list from an existing list and adds an element to the end.

    go:list_append(List, Element, NewList)
    
* `List`: a list
* `Element`: any term
* `NewList`: a variable (that will be bound to the ordered list)

## List order

Creates a new list, based on an existing list, but ordered by an order function.

    go:list_order(List, &OrderFunction, NewList)
    
* `List`: a list
* `OrderFunction`: a reference to a rule that functions as an order function
* `NewList`: a variable (that will be bound to the ordered list)

The order relation takes two entities and returns a negative number, 0, or a positive number. negative when E1 goes before E2, 0 when E1 has the same order position as E2, positive when E1 goes after E2.    
    
    by_easiness(E1, E2, R) :- if_then_else( cleartop(E1), unify(R, 1), unify(R, 0) );
    
    go:list_order(List, &by_easiness, NewList)

## List foreach

Binds each of the values of list to `Variable`, and executes `Scope` for each value.

There are two variants: one that binds a variable `ElementVar` each iteration, and one that also bindings an index (0, 1, 2, ...) 

    go:list_foreach(List, ElementVar, Scope)
    go:list_foreach(List, IndexVar, ElementVar, Scope)
    
* `List`: a list
* `IndexVar`: a variable
* `ElementVar`: a variable
* `Scope`: a relation set

## List deduplicate

Creates a new list, based on existing list, but with all duplicate elements removed. The order of the values will not change.

    go:list_deduplicate(List, NewList)
    
* `List`: a list
* `NewList`: a variable that will contain a list

## List sort

Creates a new list, based on existing list, but with all elements sorted. 

    go:list_sort(List, NewList)
    
* `List`: a list
* `NewList`: a variable that will contain a list

The function checks the types of the elements. If they are all integers, the list will be sorted from low to high. If they are all strings, they will be sorted alphabetically. Otherwise, it will cause an error.

## List length

Puts the number of elements of List in Len. 

    go:list_length(List, Len)
    
* `List`: a list
* `Len`: a variable that will contain an integer
 
## List index
 
Puts the index of the occurrence of `E` in `Index`. If there are more occurrences, more bindings will be created. 
 
     go:list_index(List, E, Index)
     
* `List`: a list
* `E`: any term
* `Index`: an variable that will contain an integer
 
## List get
 
Loads the `Index`'th term in `E` 
 
    go:list_get(List, Index, E)
     
* `List`: a list
* `Index`: an integer
* `E`: a variable that will contain the term

## List expand

Creates a new binding for each of the elements of `List`

    go:list_expand(List, E)
    
* `List`: a list
* `E`: a variable that will contain a term

## List head

Separates the head of the list from the tail.

    go:list_head(List Head, Tail)

* `List`: a list
* `Head`: a variable to contain any element
* `Tail`: a variable that will contain a list
