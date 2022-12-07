# Code

To ask the user the clarification question on which "on" they mean, you can do this

    on(A, B) :-
        go:translate('which one', Q)
        go:translate('directly on the surface', A1)
        go:translate('anywhere on top of', A2)
        go:wait_for(
            go:user_select(Q, [A1, A2], Sel)
        )
        if [Sel == 0] then
            support(B, A)
        else
            anywhere_on(A, B)
        end
    ;

The `go:wait_for` predicate waits until the predicate `go:user_select` is successful. Then the user has answered `Q` (with possible answers `A1` and `A2`, and the index of the response (0 or 1) is placed in `Sel`).

