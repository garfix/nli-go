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

Picks the first value of the variable and uses that for all bindings. Can be used when the database returns several values and one of them is enough for you.

    go:first(Name) 
    
with bindings 
    
    [{Name:'Babbage'}{Name:'Charles B.'}{Name:'Charles Babbage'}]

returns `[{Name:'Babbage'}]`
    
## exists

Checks if there currently are any bindings. The function doesn't actually do anything. It is a filler for the condition clause in a solution. This function can only be used in the condition of a solution, because this is the only relation set that is executed even with zero bindings.

    go:exists()

## other predicates

[`make_list`](list.md) also takes multiple bindings.
