## Intent

This relation is used by solutions to recognize types of problems.

    go:intent(Atom, V...)
    
* `A`: an atom
* `V`: zero or more variables    
    
An intent function has no effect; it always succeeds.    

## Back reference

The system will try to resolve E1 with the entities from the anaphora queue. It will check E1's type against the types of entities in the queue. It will also check if the value of the entity in the queue matches relation set `D`.

    go:back_reference(E1, D)
    
* `A`: a variable
* `B`: a relation set    

This allows you express "him", like this:

    go:back_reference(E1, gender(E1, male))

The recentless processed `person` entities will be processed, and the `gender()` check makes sure these persons are male.

## Definite reference

A definite reference checks not only in the anaphora queue, but in the databases as well. 

If more than one entity matches, a remark is returned to the user: "I don't understand which one you mean"

    go:definite_reference(E1, D)
    
* `A`: a variable
* `B`: a relation set
