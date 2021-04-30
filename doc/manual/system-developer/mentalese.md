# Mentalese

I called the programming language "Mentalese" because it's a mental programming language.

## Process flow

Processing is different than in Prolog. Prolog is more depth-first. Mentalese is breadth-first.

Take for example the flow

    a(X) b(X) c(X)

Prolog will do `a(X)`, finds a value, and then goes on to `b(X)` with just this value. When it fails, it will "backtrack" and find another value.

Mentalese will also do `a(X)`, find n values, and then goes on to `b(X)`. It tries all these values, and this will result in m bindings. It will take all bindings to `c(X)` and try them all out.

Why breadth-first? I didn't think about it when I made it at the time. Now I have several multi-binding functions that depend on it. So it makes some things easy that are not possible, or much harder, with depth-first, like checking how many bindings we currently have for a variable.
