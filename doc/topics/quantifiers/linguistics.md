# Linguistics

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
