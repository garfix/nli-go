# Process bindings

It seems so simple, to pass a variable binding to a relation and get a set of new bindings in return. But a surprisingly large things can go wrong if you don't know exactly what needs to be done.

An example to explain things:

    in: relationX(A, relY(B))    binding: { A:1, C:3 }
    temporary bindings: { C:'world', D:8, E:4 }
    out: [
        { A:1, C:3, B:2 },
        { A:1, C:3, B:5 }
    ]

I will now go into the things you need to think about:

## Pass only the essential bindings

In the example above, only the binding of `A` should be passed to the relation; but not the binding of `C`. The reason is that `C` may also occur as temporary variable when `relationX` is processed, and passing `C` would already give the temporary variable a value. This is always wrong, but you will only notice if `C` conflicts with a temporary variable.

The example shows that the essential variables may be found in nested relations as well (here: `B`).

## Do not return the bindings of temporary variables

While processing `relationX`, a new variable `D` may be bound. It is very important that this binding is not integrated in the result set of bindings. The result set should only contain the original bindings, appended with the new values of the variables of `relationX`.

## Create new variables for new variables

When the child relation set contains variables that have not been bound, create new variables for each of them.

## Rules

When a relation is "solved" using a rule, there are some additional things to think about.

Say our relation is `relationX(A, B)`, and it matches with `relationX(X, Y`) in a rulebase. We can do either of two things: bind the variables `X, Y` to the values of `A, B` and dive into the recursion; or: rewrite the rule using the input variables and dive into the recursion. We do the latter. This way it is easier to follow what happens to a variable in the process, as it doesn't change its name all the time.

So when the rule 

    relationX(X, Y) :- relationZ(X, Z) relationW(Z, Y)
    
matches, we first change it into 

    relationX(A, B) :- relationZ(A, Gen1) relationW(Gen1, B)
    
and only then it is processed.    

## Testing

If you add a new relation processing, these are essential things to test.
 
