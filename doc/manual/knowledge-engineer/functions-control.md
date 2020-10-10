# Control functions

## And

A variant on the boolean function `and` that works with bindings.

    and(A, B)
    
* `A`: a relation set
* `B`: a relation set    

`and(A, B)` processes `A` first. Then it processes `B`. The bindings from A are used in B.

Note: the boolean functions (`and`, `or` and `xor`) have a different meaning when used with quants. See [Nested Quants](quantification.md#nested-quants) 

## Or

A variant on the boolean function `or` that works with bindings. This is the only operator that can yield more bindings than each of its children.

    or(A, B)
    
* `A`: a relation set
* `B`: a relation set    

`or(A, B)` processes both `A` and `B`. The bindings of both are combined and doubles are removed.

## Xor

A variant on the boolean function `xor` that works with bindings. Resolves either A or B, and returns the results of the first successful one. 

    xor(A, B)
    
`xor(A, B)` processes `A` first. If it yields results, `xor` stops. Otherwise it processes `B`.     

## Not

A variant on the boolean function `not` that works with bindings. If executing `A` does not return any bindings, `not` will return its original bindings; if `A` does return bindings, `not` will not return any bindings. 

    not(A)
    
* `A`: a relation set   

## If / then / else

If `Condition` succeeds, then `Action` is executed. If not, then `Alternative` is executed. 

    if_then_else(Condition, Action, Alternative)     

## Intent

This relation is used by solutions to recognize types of problems.

    intent(Atom, V...)
    
* `A`: an atom
* `V`: zero or more variables    
    
An intent function has no effect; it always succeeds.    

## Back reference

The system will try to resolve E1 with the entities from the anaphora queue. It will check E1's type against the types of entities in the queue. It will also check if the value of the entity in the queue matches relation set `D`.

    back_reference(E1, D)
    
* `A`: a variable
* `B`: a relation set    

This allows you express "him", like this:

    back_reference(E1, gender(E1, male))

The recentless processed `person` entities will be processed, and the `gender()` check makes sure these persons are male.

## Definite reference

A definite reference checks not only in the anaphora queue, but in the databases as well. 

If more than one entity matches, a remark is returned to the user: "I don't understand which one you mean"

    definite_reference(E1, D)
    
* `A`: a variable
* `B`: a relation set

## Call

This relation just processes its single argument, that is a relation set. The purpose of this is to implement words like "tell", whose argument is a clause.

    call(S)
    
* `S`: a relation set    
