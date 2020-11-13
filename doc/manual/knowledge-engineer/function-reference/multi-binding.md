# Multi-binding predicates

These predicates take all current bindings as input, and replace these with new bindings.

These are some built-in functions from `SystemMultiBindingBase` for use in when solving problems.

## number_of

Counts the number of distinct values of `Var` in the bindings, and places the result in the `Number` value of each of the bindings. Or, if `Number` is a value, checks if this value matches the actual number of distinct values in the bindings.

    go:number_of(Var, Number)
 
* `Var`: an unbound or bound variable
* `Number`: an unbound variable or an integer

Example:

Place the number of distinct values of `E1` in `Number`

    go:number_of(E1,Number) with bindings [{E1: 5}{E1: 13}{E1: 5}]
    
returns `[{E1: 5, Number:2}{E1: 13, Number:2}{E1: 5, Number:2}]`

Check if the number of distinct values is 3
    
    go:number_of(E1,3)

If true, returns the original bindings. If false, returns an empty set. 

## first

Picks only the first binding of the results. Use `Length` to extract `Length` first bindings.

    go:first()
    go:first(Length) 
    
with bindings 
    
    [{Name:'Babbage'}{Name:'Charles B.'}{Name:'Charles Babbage'}]

returns `[{Name:'Babbage'}]`

## last

Picks only the last binding. Use `Length` to extract `Length` last bindings.

    go:last()
    go:last(Length) 
    
with bindings 
    
    [{Name:'Babbage'}{Name:'Charles B.'}{Name:'Charles Babbage'}]

returns `[{Name:'Charles Babbage'}]`

## get

Picks the `N-1`th binding. Use `Length` to extract `Length` bindings, starting with the `N-1`th.

    go:get()
    go:get(Length) 
    
with bindings 
    
    [{Name:'Babbage'}{Name:'Charles B.'}{Name:'Charles Babbage'}]

returns `[{Name:'Charles Babbage'}]`

## largest 

Takes in all current bindings and removes the ones whose E value is not the largest number.

    go:largest(E) 
    
## smallest 

Takes in all current bindings and removes the ones whose E value is not the smallest number.

    go:smallest(E) 
    
## order 

Sorts the bindings by the values of the variable `E`, either ascending (`asc`) or descending (`desc`):

    go:order(E, asc)    
    go:order(E, desc)
    
Values of E must be either all numbers, or all strings.    

## exists

Checks if there currently are any bindings. The function doesn't actually do anything. It is a filler for the condition clause in a solution. This function can only be used in the condition of a solution, because this is the only relation set that is executed even with zero bindings.

    go:exists()

## other predicates

[`make_list`](list.md) also takes multiple bindings.
