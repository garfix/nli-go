# Keywords

Mentalese has some keywords that make programs easier to read:

## if then (else) end

This is the common if then construction:

    if go:not(cleartop(E2)) do_find_free_space(E2, E1, X1, Y1) then
        do_put_on_position(ParentEventId, E1, E2, X1, Y1)
    end

or if/then/else

    if go:equals(Sel, 0) then
        support(B, A)
    else
        anywhere_on(A, B)
    end

The whitespacing is not required, but this is the preferred way of writing

## break

`break` breaks a loop and keeps the bindings built so far.

## cancel

`cancel` breaks a loop and discards all bindings built so far.

## return

`return` ends a procedure immediately, succeeding (with bindings)

## fail

`fail` ends a procedure immediately, failing (no bindings)
