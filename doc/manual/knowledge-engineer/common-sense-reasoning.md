# Common sense reasoning

NLI-GO provides rules that allow you to solve a problem in a way that is natural for a human.

The rules take the form of 

    predX(A, B) :- predY(A, C) predZ(C, B);
    
Where `predX(A, B)`, `predY(A, C)` and `predZ(C, B)` are relations (formally: atoms).    
    
## Interpretation    
    
This rule has two useful interpretations:

1) I believe `predX(A, B)` to be true if both `predY(A, C)` and `predZ(C, B)` can be shown to be true.
2) To achieve the goal `predX(A, B)`, both `predY(A, C)` and `predZ(C, B)` need to succeed.

This is the dual nature of atoms: they can be regarded as either a goal to be attained, or a statement to be proven. In order to make the interpretation more explicit, you can start a goal-atom with the prefix "do_", as in "do_stack_up".

When an atom such as `predX(A, B)` is executed, the result is not a truth value (`true`, `false`). This is because the system doesn't have access to the world directly, it merely has knowledge, or even just beliefs, about the world.   
  
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
