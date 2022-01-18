# Procedures as functions

NLI-GO is based on procedures like `process(A, B)`. Procedures take in values and return values via their arguments.

If you want to use the result of a procedure to use it somewhere else, you need to assign it first.

    add(A, B, C)
    process(C)

Here `C` is a temporary variable that serves to other purpose than to hold the return value of `add`. Can't we just use functions like this?

    process(add(A, B))

Yes we can, except that the function `add(A, B, C)` has three arguments and we need to specify which contains the return value. Do it like this:

    process(add(A, B, rv))

Internally, NLI-GO replaces `rv` by a temporary variable, whose value replaces the call `add(A, B, rv)` in the call to `process`.

Note: in the current implementation functions return only a single binding for a value; the first one returned.
