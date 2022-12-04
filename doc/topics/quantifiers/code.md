# Code

## Quantifier

In order to declare the quantifier we use the relation `quantifier`, specify variables for the two numbers, and give the relation set for the quantification. 

    every:              quantifier(Result, Range, equals(Result, Range))
    more than two:      quantifier(Result, Range, greater_than(Result, 2))
    two or three:       quantifier(Result, Range, or(P1, equals(Result, 2), equals(Result, 3)))

## Quants 

Whenever an NP is used, it is quantified: its sense is a `quant()`. The prototypical case is:

    { rule: np(R1) -> qp(_) nbar(R1),                                      sense: quant($qp, R1, $nbar) }
    
Here `$nbar` means: include the sense of the second right-hand structure (the `nbar`) in this position, and `$qp` means: include the sense of the first structure. This `$qp` must be a `quantifier()`. 
    
These are the arguments of the `quant`:

    quant(
        Quantifier,
        RangeVariable,
        RangeRelations        
    )
    
and these are the arguments of the the `quantifier`:

    quantifier(
        ResultCountVariable,
        RangeCountVariable,
        ScopeSet
    )    

Here is the typical case for the `qp`. Note that the `quantifier` relation is formed only once, in the rewrite of `qp` to `quantifier`. The sense of the quantifier can be simple (see 'every') or compound ('or').    

    { rule: qp(_) -> quantifier(Result, Range),                                                         sense: quantifier(Result, Range, $quantifier) }
    { rule: quantifier(Result, Range) -> 'every',                                                       sense: equals(Result, Range) }
    { rule: quantifier(Result, Range) -> 'some',                                                        sense: greater_than(Result, 0) }
    { rule: quantifier(Result, Range) -> number(N1),                                                    sense: equals(Result, N1) }
	{ rule: quantifier(Result, Range) -> quantifier(Result, Range) 'or' quantifier(Result, Range),	    sense: or(P1, $quantifier1, $quantifier2) }

## Quant check

The quant is only useful when combined with a parent relation (typically a verb). You need to specify explicity that the quant is used. If there is more than one quant, the order of the quants can be given. An example:

    { rule: np_comp4(P1) -> np(E1) marry(P1) 'to' np(E2),                    sense: check($np1, check($np2, marry(P1, E1, E2))) }
    
Imagine the sentence: "Did all these men marry two women?". Resolving this question means going through all the men, one-by-one, and for each of them counting the women that were married to them. If one of them married only one woman, the answer is no.     
    
`check` says: apply the quants from the right-hand positions 1 (`$np1`), which is the sense of the `np(E1)` and 4 (`$np2`) and use them in that order. When the sense is built, the result looks like this:

    check(
        quant(Q1, E1, ...), 
        marry(P1, E1, E2)
    )     

`check` has a set of quantifiers, and a _scope_ that consists of zero or more relations (`marry(P1, E1, E2)`).

In this example the quant for E1 precedes that of E2, but the order does not need to match the order of the variables in the scope.

It is important to understand the way `check` is evaluated. I will sketch the process here briefly. Note that the quants are nested, and that the inner loop uses a single value from the range of the outer quant, and goes through all range values of the inner loop.

    b1 = []binding
    foreach E1-range-set as E1 in outer range {
        b2 = []binding
        foreach E2-range-set as E2 in inner range {
            execute scope, bound with single E1 and E2, and add binding to b2
        }
        check quantifier Q2 with disinct-values(E2 in b2), may fail (b2 = [])
        add b2 to b1
    }
    check quantifier Q1 with disinct-values(E1 in b1), may fail (b1 = [])
    return b1
    
This is the process for 2 quants. The number of quants is often 1, and possibly more than 2.    

## Quant foreach

The function `do` is exactly like find, with one important distinction: `do` checks the quantifier _during_ the loop as well.

    foreach E2-range-set as E2 in inner range {
        execute scope, bound with single E1 and E2, and add binding to b2
        check quantifier Q2 with disinct-values(E2 in b2), breaks on success
    }  

This relation is needed for different kinds of relations: imperative ones, like `pick_up()`.

Imagine now this sentence: "pick up two blocks". 

Handling this question with `check` amounts to picking up all blocks and then checking if there were two that were picked up. This is clearly nonsense. `do` goes through all blocks, and attempts to pick them up. As soon as the quantifier `2` matches, it stops.

The difference between `check` and `do` is that `do` stops when it has enough, while `check` continues. Use `check` with interrogative relations and `do` with imperative relations. 

## Unquantified nouns

Unquantified nouns are nouns that are not preceded by a determiner or quantifier. An example is 

    blocks 

We still use a quantifier in this situation, because it is easier to treat all NP's as quants. But the the quantifier is 

    none
    
An example grammar rule is

    { rule: np(E1) -> nbar(E1),                                            sense: quant(none, E1, $nbar) }       

The system will find as much entities that match the scope as it can, and it always succeeds.

## Nested quants

To model a compound NP like "both of the red blocks and either a green cube or a pyramid" You can nest quants with boolean operators `and`, `or` and `xor`. For example

    { rule: np(E1) -> 'either' np(E1) 'or' np(E1),                         sense: or(_, $np1, $np2) }
        
The meaning of the operators corresponds with what you might expect, but here's a more detailed description:

`xor` means: "either A or B, but not both". First the range of A is determined and used to evaluate the scope. Only if this produces no bindings, the range of B is determined and used.

`or` means: "A or B, or both". First the range of A is determined and used to evaluate the scope. Then the range of B. Then the results are combined.

`and` means: "A and B must match". First the range of A is determined and used to evaluate the scope. Only if this produces results the range of B is determined and evaluated.     
