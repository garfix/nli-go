# Variables

Variables in NLI-GO are write-once. They may be introduced at any time without a declaration of type. Once they have a value, they keep it for the rest of the flow. (But see rule scope for exception)

## Scope

In a sequence of relations

    a(X, Y) b(Y, Z) d(Z, R)
    
each of the relations may be an invocation of a rule which itself instantiates many variables. None of these variables will bubble up to this level.

## Rule scope

When a rule is instantiated, a local scope is created for it. Within this scope it is possible to create rewritable variables, with the `go:let()` function:

    go:let(X, 123)
    
This rewritable variable may ge assigned a new variable, just like this

    go:let(X, 456)
    
