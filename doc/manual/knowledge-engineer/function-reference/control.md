# Control predicates

These predicates just affect the flow of data in an application.

## And

A variant on the boolean function `and` that works with bindings.

    go:and(A, B)
    
* `A`: a relation set
* `B`: a relation set    

`go:and(A, B)` processes `A` first. Then it processes `B`. The bindings from A are used in B.

Note: the boolean functions (`and`, `or` and `xor`) have a different meaning when used with quants. See [Nested Quants](quantification.md#nested-quants) 

## Or

A variant on the boolean function `or` that works with bindings. This is the only operator that can yield more bindings than each of its children.

    go:or(A, B)
    
* `A`: a relation set
* `B`: a relation set    

`go:or(A, B)` processes both `A` and `B`. The bindings of both are combined and doubles are removed.

## Xor

A variant on the boolean function `xor` that works with bindings. Resolves either A or B, and returns the results of the first successful one. 

    go:xor(A, B)
    
`go:xor(A, B)` processes `A` first. If it yields results, `xor` stops. Otherwise it processes `B`.     

## Not

A variant on the boolean function `not` that works with bindings. If executing `A` does not return any bindings, `not` will return its original bindings; if `A` does return bindings, `not` will not return any bindings. 

    go:not(A)
    
* `A`: a relation set   

## Call

This relation just processes its single argument, that is a relation set. The purpose of this is to implement words like "tell", whose argument is a clause.

    go:call(S)
    
* `S`: a relation set

## Ignore

This relation just processes its single argument, that is a relation set. The difference with `call` is that `ignore` always succeeds. 

    go:ignore(S)

* `S`: a relation set

## Range foreach

For loop over an integer range Start .. End. `Scope` will be called with Variable instantiated to each of the numbers in [Start..End], including Start and End. 

    go:range_foreach(Start, End, Variable, Scope)
    
* `Start`: an integer
* `End`: an integer
* `Variable`: a variable    
* `Scope`: a relation set    

Example:

    go:range_foreach(1, 10, I,
        go:multiply(Result, I, Result)
    )     

## Wait for

This relation tries `Condition` until it succeeds.

Under the hood, it tries `Condition` once. If it succeeds, `wait_for` succeeds. If it fails, `Condition` is restacked and the process ends without breaking up the stack. Next time the process is executed, the `Condition` will be tried again.

    go:wait_for(Condition)

* `Condition`: a relation set

Example:

    go:wait_for(
        go:which_one(['George', 'Jack', 'Bob'], SelectionIndex)
    )

## Slot

Gets or sets a process slot value. 

Slots are like global variables, in the space of a process. They should be used only for a very limited set of slots, like `locale`.  

    go:slot(locale, 'nl_NL')

    go:slot(locale, Locale)

## Isa

Checks an inheritance relation. The sort of the variable is deduced from the relations in the input sentence during the relationization process. 

    go:isa(Var, Sort)

* `Var`: a variable
* `Sort`: a sort atom 

Example

    go:isa(E1, block)


## Has Sort

To declare the sort of an entity:

    go:has_sort(Id, Sort)

* `Id`: an id
* `Sort`: an atom

To declare the supersort of a sort

  go:has_sort(Sort, SuperSort)

* `Sort`: an id
* `SuperSort`: an atom

To query the sort of an entity:

  go:has_sort(Id, Sort)

* `Id`: an id / a variable bound to an id
* `Sort`: a variable

Example

    go:has_sort(E, Sort)
    go:has_sort(E, block)
    go:has_sort(block, object)

