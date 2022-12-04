# Considerations

## Strong Negation

Atoms can also take a negative form
    
    `-predH(A, B)` :- predI(A, C) -predJ(C, B);
             
The interpretation of this rule is: 

1) I believe `predH(A, B)` to be false if I believe `predI(A, C)` to be true and `predJ(C, B)` not to be true

Example interpretations of `-predH` are "it does not rain", "is not red", "is not on". This kind of negation is different from `not` above. Whereas `not` inverts `known`/`unknown`, `-` inverts the meaning of the predicate itself. `-raining()` means "dry or foggy or whatever, but not raining". It means "everything else". `-red` means "blue, yellow, green, purple, etc, etc" It is an affirmation of the positive belief to the complement of the original predicate: "I believe `-predH()` to be true". 

In a Closed World `-` and `not` come down to the same thing, since it assumes that what is unknown is also untrue. NLI-GO takes a broader view in order to deal with the full power of natural language. It takes the Open World view. 

However, when Open World proves to be too unrestrictive for a use case, you can make exceptions, and say: "for this type of predicates I need to closed world". Here's an example:

    `-predH(A, B)` :- not(predH(A, B));
    
This means simply: I believe `predH(A, B)` to be false if I have no knowledge about `predH(A, B)`. Examples are abound: if I can't find a reservation for customer C, I believe he has not made a reservation. (Full open world would have leave you in doubt as to wether the reservation had been made.)
     
More about this form of negation, called "strong negation" can be read [on Wikipedia](https://en.wikipedia.org/wiki/Stable_model_semantics#Strong_negation) 

So when do you use `not` and when `-`?

* Use `not` when you mean "failed" (if it's a goal) or "not found"/"not proven" (if it's a statement)
* Use `-` when you mean "know not to be the case"; this operator is not applied to goals

