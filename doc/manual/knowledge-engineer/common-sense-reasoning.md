# Common sense reasoning

NLI-GO provides rules that allow you to solve a problem in a way that is natural for a human.

The rules that solve goals take the form of 

    predX(A, B) :- predY(A, C) predZ(C, B);
    
## Interpretation    
    
This rule has two useful interpretations

1) I believe `predX(A, B)` to be true if both `predY(A, C)` and `predZ(C, B)` can be shown to be true. This
2) In order to reach `predX(A, B)` both `predY(A, C)` and `predZ(C, B)` need to succeed.

When an atom such as `predX(A, B)` is executed, the result is not a truth value (`true`, `false`). This is because the system doesn't have access to the world directly, it merely has knowledge, or even merely beliefs, about the world.   
  
 Rather than truth values, the rules produce a set of variable bindings, like this:
    
    [
        {A:1, B:2}
        {A:3, B:2}
        {A:5, B:8}
        {A:9, B:4}
    ]
         
An empty set of bindings corresponds to `known`/`success` (i.e. bound in the set of beliefs, giving results) and a non-empty set corresponds to `unknown`/`failure` (not bound in the set of beliefs, giving no results).

## Composition

Three logical predicates act on these relations to produce new bindings:

* `and(P1(), P2())` evaluates `P1()` and `P2()` and returns their bindings.
* `or(P1(), P2())` evaluates `P1()` and if it has bindings, it returns these; if not, it returns the bindings of `P2()`
* `not(P1())` evaluates `P1()`; if it has bindings, `not()` returns no bindings; if it has no bindings, `not()` returns its original bindings

## Negation

The goals can also take a negative form
    
    `-predH(A, B)` :- predI(A, C) -predJ(C, B);
             
The interpretation of this rule is: 

1) I believe `predH(A, B)` to be false if I believe `predI(A, C)` to be true and `predJ(C, B)` not to be true

Example interpretations of `-predH` are "it does not rain", "is not red", "is not on". This kind of negation is different from `not` above. Whereas `not` inverts `known`/`unknown`, `-` inverts the meaning of the predicate itself. `-raining()` means "dry or foggy or whatever, but not raining". It means "everything else". `-red` means "blue, yellow, green, purple, etc, etc" It is an affirmation of the positive belief to the complement of the original predicate: "I believe `-predH()` to be true". 

In a Closed World `-` and `not` come down to the same thing, since it assumes that what is unknown is also untrue. NLI-GO takes a broader view in order to deal with the full power of natural language. It takes the Open World view. 

However, when Open World proves to be too unrestrictive for a use case, you can make exceptions, and say: "for this type of predicates I need to closed world". Here's an example:

    `-predH(A, B)` :- not(predH(A, B));
    
This means simply: I believe `predH(A, B)` to be false if I have no knowledge about `predH(A, B)`. Examples are abound: if I can't find a reservation for customer C, I believe he has not made a reservation. (Full open world would have leave you in doubt as to wether the reservation had been made.)
     
More about this form of negation, called "strong negation" can be read [on Wikipedia](https://en.wikipedia.org/wiki/Stable_model_semantics#Strong_negation) 

 