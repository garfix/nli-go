# Quantification

Quantification is a central concept in the modelling of meaning. Each NP is quantified, which means that there are an exact number of entities involved in each NP. In "3 children" this number is 3. In "Some children" this number is higher than 0. The same holds for "children". In "All children" the number equals the number of children that results from the rest of the predication.     

## Generalized quantifiers

Traditional predicate logic just has the quantifiers (quantors) `exists` and `all`.

Natural language also allows quantifiers like "more than 2", "two or three", "between 5 and 10" and "most". In order to understand how these quantifiers are modelled, we must go into some theory.

Let's take the sentence:

    Did all red balls fall from the table?
    
When is this sentence true? In order to answer this we distinguish between "all red balls" (the _quantification_) and "fall from the table" (the _scope_). In order to determine if the sentence is true, we must go through all entities E that are "red balls", and we will find a set of identifiers for these red balls. This set of (unique) identifiers we call the _range_. The range may have 3 values (let's call this `Range`). Now we try the scope for each of these values. Ball `1`, did it fall from the table? Yes. Ball `2`, did it fall from the table? No. Ball `3`? Yes, it did. We found that 2 red balls actually fell from the table (call it `Result`). Now we can answer the question, as _all_ simply means: `equals(Result, Range)`. This returns the empty set, so that means: No.

"all" is actually one of the few quantifiers that requires two numbers. Most quantifiers just use a single value (`Result`). They don't care how many entities where in the range, as long as a specific number made it into the result. An example is:

    Did at least two red balls fall from the table?
    
Again all red balls are checked and the number of red balls that actually fell from the table is counted in `Result`. The answer to the question is now `greater_than(Result, 2)`    

In order to declare the quantifier we use the relation `quantifier`, specify variables for the two numbers, and give the relation set for the quantification. 

    every:              quantifier(Result, Range, equals(Result, Range))
    more than two:      quantifier(Result, Range, greater_than(Result, 2))
    two or three:       quantifier(Result, Range, or(P1, equals(Result, 2), equals(Result, 3)))

## Quants 

Whenever an NP is used, it is quantified: its sense is a `quant()`. The prototypical case is:

    { rule: np(R1) -> qp(_) nbar(R1),                                      sense: quant(sem(1), R1, sem(2)) }
    
Here `sem(2)` means: include the sense of the second right-hand structure (the `nbar`) in this position, and `sem(1)` means: include the sense of the first structure. This `sem(1)` must be a `quantifier()`. 
    
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

    { rule: qp(_) -> quantifier(Result, Range),                                                         sense: quantifier(Result, Range, sem(1)) }
    { rule: quantifier(Result, Range) -> 'every',                                                       sense: equals(Result, Range) }
    { rule: quantifier(Result, Range) -> 'some',                                                        sense: greater_than(Result, 0) }
    { rule: quantifier(Result, Range) -> number(N1),                                                    sense: equals(Result, N1) }
	{ rule: quantifier(Result, Range) -> quantifier(Result, Range) 'or' quantifier(Result, Range),	    sense: or(P1, sem(1), sem(3)) }

## Find

The quant is only useful when combined with a parent relation (typically a verb). You need to specify explicity that the quant is used. If there is more than one quant, the order of the quants can be given. An example:

    { rule: np_comp4(P1) -> np(E1) marry(P1) 'to' np(E2),                    sense: quant_check([sem(1) sem(4)], marry(P1, E1, E2)) }
    
Imagine the sentence: "Did all these men marry two women?". Resolving this question means going through all the men, one-by-one, and for each of them counting the women that were married to them. If one of them married only one woman, the answer is no.     
    
`find` says: apply the quants from the right-hand positions 1 (`sem(1)`), which is the sense of the `np(E1)` and 4 (`sem(4)`) and use them in that order. When the sense is built, the result looks like this:

    quant_check(
        [
            quant(Q1, E1, [...]) 
            quant(Q2, E2, [...])
        ], [
            marry(P1, E1, E2)
        ]
    )     

`find` has a set of quantifiers, and a _scope_ that consists of zero or more relations (`marry(P1, E1, E2)`).

In this example the quant for E1 precedes that of E2, but the order does not need to match the order of the variables in the scope.

It is important to understand the way `find` is evaluated. I will sketch the process here briefly. Note that the quants are nested, and that the inner loop uses a single value from the range of the outer quant, and goes through all range values of the inner loop.

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

## Do

The function `do` is exactly like find, with one important distinction: `do` checks the quantifier _during_ the loop as well.

    foreach E2-range-set as E2 in inner range {
        execute scope, bound with single E1 and E2, and add binding to b2
        check quantifier Q2 with disinct-values(E2 in b2), breaks on success
    }  

This relation is needed for different kinds of relations: imperative ones, like `pick_up()`.

Imagine now this sentence: "pick up two blocks". 

Handling this question with `find` amounts to picking up all blocks and then checking if there were two that were picked up. This is clearly nonsense. `do` goes through all blocks, and attempts to pick them up. As soon as the quantifier `2` matches, it stops.

The difference between `find` and `do` is that `do` stops when it has enough, while `find` continues. Use `find` with interrogative relations and `do` with imperative relations. 

## Unquantified nouns

Unquantified nouns are nouns that are not preceded by a determiner or quantifier. An example is 

    blocks 

We still use a quantifier in this situation, because it is easier to treat all NP's as quants. But the the quantifier is 

    none
    
An example grammar rule is

    { rule: np(E1) -> nbar(E1),                                            sense: quant(none, E1, sem(1)) }       

The system will find as much entities that match the scope as it can, and it always succeeds.

## Nested quants

To model a compound NP like "both of the red blocks and either a green cube or a pyramid" You can nest quants with boolean operators `and`, `or` and `xor`. For example

    { rule: np(E1) -> 'either' np(E1) 'or' np(E1),                         sense: xor(_, sem(2), sem(4)) }
        
The meaning of the operators corresponds with what you might expect, but here's a more detailed description:

`xor` means: "either A or B, but not both". First the range of A is determined and used to evaluate the scope. Only if this produces no bindings, the range of B is determined and used.

`or` means: "A or B, or both". First the range of A is determined and used to evaluate the scope. Then the range of B. Then the results are combined.

`and` means: "A and B must match". First the range of A is determined and used to evaluate the scope. Only if this produces results the range of B is determined and evaluated.     
