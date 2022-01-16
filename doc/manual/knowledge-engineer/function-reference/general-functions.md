# Functions

These are some built-in functions from `SystemFunctionBase` to be used everywhere where relations are used.

Each function takes in a single binding, and returns either an extended binding when it succeeds, or no binding when it fails.  

When an argument has a specific type, let's say a string, this means that it can be a string, or a variable bound to a string value.

## split

Splits a string `Whole` into parts using a `separator`. The result is placed in one or more variables `Part1`, `Part2`, etc.

    go:split(Whole, Separator, Part1, Part2, ...)
 
* `Whole`: a string
* `Separator`: a string
* `Part1`, `Part2`: free variables (for strings)

Example:

    go:split(Fullname, " ", Firstname, Lastname) 

## join

Joins two or more strings `Part1`, `Part2` ... into a `Whole`, using a `separator`

    go:join(Whole, Separator, Part1, Part2, ...)
 
* `Whole`: a variable (for a string)
* `Separator`: a string
* `Part1`, `Part2`: strings

Example:

    go:join(Fullname, " ", Firstname, Insertion, Lastname)
    
## concat

Same as `go:join`, with empty string as separator.

    go:concat(Whole, Part1, Part2, ...)

* `Whole`: a variable (for a string)
* `Part1`, `Part2`: strings
   
# unify
    
This function unifies two terms. This can be used for assignment and destructuring. Assignment only works for variables that had not been previously assigned.

    go:unify(T1, T2)
    
* `T1`: a free variable, or any other term
* `T2`: a free variable, or any other term

Examples:

If N1 is an unbound variable, this function is an assignment (N1 becomes 2).
If N1 is a bound variable, this function checks if both the type and the value are identical.

    go:unify(N1, 2)
    
Destructuring. If Q1 holds a quant() relation, this equals, binds its arguments to new variables (`R1`).
    
    go:unify(Q1, quant(_, _, R1, _)


## add

Adds two numbers `N1` and `N2` and places the result in `Sum`. If `Sum  is a number and it is not the sum of the arguments, the function fails.  

    go:add(N1, N2, Sum)
    
* `N1`: a number
* `N2`: a number
* `Sum`: a variable (to contain a number) or a number

## subtract

Subtracts two numbers `N1` and `N2` and places the result in `Diff`. If `Diff  is a number and it is not the diff of the arguments, the function fails.  

    go:subtract(N1, N2, Diff)
    
* `N1`: a number
* `N2`: a number
* `Sum`: a variable (to contain a number) or a number

## multiply

Multiplies two numbers `N1` and `N2` and places the result in `Product`. If `Product  is a number and it is not the diff of the arguments, the function fails.  

    go:multiply(N1, N2, Product)
    
* `N1`: a number
* `N2`: a number
* `Product`: a variable (to contain a number) or a number

## divide

Divides two numbers `N1` and `N2` and places the result in `Result`.

    go:divide(N1, N2, Result)

* `N1`: a number
* `N2`: a number
* `Result`: a variable (to contain a number) or a number

## min

Sets `Min` to the smallest of `N1` and `N2.  

    go:min(N1, N2, Min)
    
* `N1`: a number
* `N2`: a number
* `Sum`: a variable (to contain a number) or a number

## compare

Adds two strings `N1` and `N2` and places the result in `R`. If N1 < N2, then R = -1; if N1 = N2, them R = 0; if N1 > N2, then R = 1;  

    go:compare(N1, N2, R)
    
* `N1`: a string
* `N2`: a string
* `R`: a variable (to contain a number)

This function is useful in order functions.

## date_today

Places the date of today in the variable in the form YYYY-mm-dd. If `D1` is not free, and does not contain today's date, the function fails.

    go:date_today(D1)

* `D1`: a variable, to contain a date

## date_subtract_years

Age calculation. If `D1` and `D2` contain dates, `Years` will be assigned the difference between these dates, in years (rounded down).  

    go:date_subtract_years(D1, D2, Years)
    
* `D1`: a variable, to contain a date
* `D2`: a variable, to contain a date
* `Years`: a variable, to contain a date

## log

Prints `Str` for debugging purposes.

    go:log(Str)
    go:log(Str1, Str2, ...)
    
* `Str`: a string value

## uuid

Generates a random 16 digit hexadecimal number.

    go:uuid(V)

* `V`: a variable
