# Functions for generation

Functions used in the preparation phase of a solution, to help build the response relation set. 

## Make and

Creates a nested and-structure based on all values of `E1` in the binding.

    go:make_and(E1, And)
    
* `E1`: a variable holding an id
* `And`: a relation set consisting of `and()`s

`E1` is an input parameters. `And` is output.
