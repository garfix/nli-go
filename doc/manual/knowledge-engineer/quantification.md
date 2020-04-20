# Quantification

Quantification is a central concept in the modelling of meaning. Each NP has a quantity, like "some", "all", "at least 2", but also implicit as with proper nouns (quantity = 1). While the default quantifier for other entities is `exists` / `some`, this is not so for NP's. NP's are always explicitly quantified. NP's typically serve as the arguments for verbs, but they may be used in other relations as well.

## Quants 

It is customary to use the part-of-speech `np` only for quantified entities. Whenever an NP is used, it is quantified: its sense is a `quant()`. The prototypical case is:

    { rule: np(R1) -> qp(Q1) nbar(R1),                                      sense: quant(Q1, sem(1), R1, sem(2)) }
    
Here `sem(2)` means: include the sense of the second right-hand structure (the `nbar`) in this position. Thus `quant` is a second-order relation that nests the senses of its dependent phrases.    
    
These are the parts:

    quant(
        QuantifierVariable,
        QuantifierRelations,
        RangeVariable,
        RangeRelations        
    )

The quant consists of a _quantifier_ and a _range_. 

The quantifier consists of a variable and a relation set and specifies the quantity of the NP ("all", "2 or more", etc). The variable is only needed if the relation set contains a variable. In most cases like when the relation set is `all(E1)` or `number(N1)` the variable can be the anonymous variable `_`.  

The range also consists of a variable and a relation set, and describes the entities involved ("child", "country", "block that i asked you to pick up"). Since the relation set often contains multiple variables, the range variable is needed to specify the entity that matters.

There are four built-in quantifiers:

* the(_) : requires that the range consists of 1 entity
* all(_) : specifies that all entities in the range are needed
* some(_) : specifies that at least one entity in the range is needed
* number(N) : specificies that N entities are needed

Quantifiers like "at least two" are not yet possible.

## Find

The quant is only useful when combined with a parent relation (typically a verb). You need to specify explicity that the quant is used. If there is more than one quant, the order of the quants can be given. An example:

    { rule: np_comp4(P1) -> np(E1) marry(P1) 'to' np(E2),                    sense: find([sem(1) sem(4)], marry(P1, E1, E2)) }
    
Imagine the sentence: "Did all these men marry two women?". Resolving this question means going through all the men, one-by-one, and for each of them counting the women that were married to them. If one of them married only one woman, the answer is no.     
    
`find` says: apply the quants from the right-hand positions 1 (`sem(1)`), which is the sense of the `np(E1)` and 4 (`sem(4)`) and use them in that order. When the sense is built, the result looks like this:

    find(
        [
            quant(Q1, [...]], E1, [...]) 
            quant(Q2, [...], E2, [...])
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

The difference between `find` and `do` is that `do` stops when it has enough, while `find` contines. Use `find` with interrogative relations and `do` with imperative relations. 
