# Nested functions

## And

A variant on the boolean function `and` that works with bindings.

    and(P1, A, B)
    
* `A`: a relation set
* `B`: a relation set    

`and(C1, A, B)` processes `A` first. Then it processes `B`. The bindings from A are used in B.

## Or

A variant on the boolean function `or` that works with bindings.

    or(P1, A, B)
    
* `A`: a relation set
* `B`: a relation set    

`or(C1, A, B)` processes `A` first. If it yields results, `or` stops. Otherwise it processes `B`.

## Not

A variant on the boolean function `not` that works with bindings. If executing `A` does not return any bindings, `not` will return its original bindings; if `A` does return bindings, `not` will not return any bindings. 

    not(A)
    
* `A`: a relation set   

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

## Do / Find

Perform the relation set `S` while iterating over the entities described by `Q`.

    do(Q, S)
    find(Q, S)
    
* `Q`: a relation set of `quant`s.
* `S`: a relation set    

Check [quantification](quantification.md) for more information.
