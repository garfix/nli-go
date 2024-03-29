# Keywords

Mentalese has some keywords that make programs easier to read:

## assignment

To assign a value to a free variable, do

    [X := n]

Where `X` is a (unmutable) variable and `n` any term. 

The same expression is used for mutable variables

    [:X := n]

any existing value is overwritten by the new value. 

## equals

    [T1 == T2]

This expression compares two terms.

* `T1`: a free variable, or any other term
* `T2`: a free variable, or any other term

## not_equals

This expression just compares two terms. If either their types or their values are unequal, it fails

    [T1 != T2]

* `T1`: a free variable, or any other term
* `T2`: a free variable, or any other term


## greater_than

Compares two integers. Succeeds if N1 > N2.

    [N1 > N2]

* `N1`: a string, representing an integer
* `N2`: a string, representing an integer

The function does not bind new variables. It just removes existing bindings if the comparison fails.

## less_than

Compares two integers. Succeeds if N1 < N2.

    [N1 < N2]

* `N1`: a string, representing an integer
* `N2`: a string, representing an integer

The function does not bind new variables. It just removes existing bindings if the comparison fails.

## greater_than_equals

Compares two integers. Succeeds if N1 >= N2.

    [N1 >= N2]

* `N1`: a string, representing an integer
* `N2`: a string, representing an integer

The function does not bind new variables. It just removes existing bindings if the comparison fails.

## less_than_equals

Compares two integers. Succeeds if N1 <= N2.

    [N1 <= N2]

* `N1`: a string, representing an integer
* `N2`: a string, representing an integer

The function does not bind new variables. It just removes existing bindings if the comparison fails.

## add

Adds two numbers `N1` and `N2`. The expression is then replaced by the result, and can be used in another expression. 

    [N1 + N2]

* `N1`: a number
* `N2`: a number

## subtract

Subtracts two numbers `N1` and `N2`. The expression is then replaced by the result, and can be used in another expression. 

    [N1 - N2]

* `N1`: a number
* `N2`: a number

## multiply

Multiplies two numbers `N1` and `N2`. The expression is then replaced by the result, and can be used in another expression.

    [N1 * N2]

* `N1`: a number
* `N2`: a number

## divide

Divides two numbers `N1` and `N2`. The expression is then replaced by the result, and can be used in another expression.

    [N1 / N2]

* `N1`: a number
* `N2`: a number

## if then (else) end

This is the common if then construction:

    if go:not(cleartop(E2)) do_find_free_space(E2, E1, X1, Y1) then
        do_put_on_position(ParentEventId, E1, E2, X1, Y1)
    end

or if/then/else

    if [Sel == 0] then
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

`fail` ends a scope immediately, failing (no bindings)
