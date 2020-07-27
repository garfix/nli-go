# Functions

These are some built-in functions from `SystemFunctionBase` to be used everywhere a match is made or a solution is found.

Each function takes in a single binding, and returns either an extended binding when it succeeds, or no binding when it fails.  

When an argument has a specific type, let's say a string, this means that it can be a string, or a variable bound to a string value.

## split

Splits a string `Whole` into parts using a `separator`. The result is placed in one or more variables `Part1`, `Part2`, etc.

    split(Whole, Separator, Part1, Part2, ...)
 
* `Whole`: a string
* `Separator`: a string
* `Part1`, `Part2`: free variables (for strings)

Example:

    split(Fullname, " ", Firstname, Lastname) 

## join

Joins two or more strings `Part1`, `Part2` ... into a `Whole`, using a `separator`

    join(Whole, Separator, Part1, Part2, ...)
 
* `Whole`: a variable (for a string)
* `Separator`: a string
* `Part1`, `Part2`: strings

Example:

    join(Fullname, " ", Firstname, Insertion, Lastname)
    
## concat

Same as `join`, with empty string as separator.

    concat(Whole, Part1, Part2, ...)

* `Whole`: a variable (for a string)
* `Part1`, `Part2`: strings

## greater_than

Compares two integers. Succeeds if N1 > N2.

    greater_than(N1, N2)
    
* `N1`: a string, representing an integer
* `N2`: a string, representing an integer

The function does not bind new variables. It just removes existing bindings if the comparison fails.

## less_than

Compares two integers. Succeeds if N1 < N2.

    less_than(N1, N2)
    
* `N1`: a string, representing an integer
* `N2`: a string, representing an integer

The function does not bind new variables. It just removes existing bindings if the comparison fails.

## greater_than_equals

Compares two integers. Succeeds if N1 >= N2.

    greater_than_equals(N1, N2)
    
* `N1`: a string, representing an integer
* `N2`: a string, representing an integer

The function does not bind new variables. It just removes existing bindings if the comparison fails.

## less_than_equals

Compares two integers. Succeeds if N1 <= N2.

    less_than_equals(N1, N2)
    
* `N1`: a string, representing an integer
* `N2`: a string, representing an integer

The function does not bind new variables. It just removes existing bindings if the comparison fails.

## equals

This function compares two terms. Next to its obvious comparison function, it is also powerful as a destructuring function.

    equals(T1, T2)
    
* `T1`: a free variable, or any other term
* `T2`: a free variable, or any other term

Examples:

If N1 is an unbound variable, this function is an assignment (N1 becomes 2).
If N1 is a bound variable, this function checks if both the type and the value are identical.

    equals(N1, 2)
    
Destructuring. If Q1 holds a quant() relation, this equals, binds its arguments to new variables (`R1`).
    
    equals(Q1, quant(_, _, R1, _)        
    
# unify
    
This function unifies two terms. This can be used for assignment and destructuring. Assignment only works for variables that had not been previously assigned.

    unify(T1, T2)
    
* `T1`: a free variable, or any other term
* `T2`: a free variable, or any other term

Examples:

If N1 is an unbound variable, this function is an assignment (N1 becomes 2).
If N1 is a bound variable, this function checks if both the type and the value are identical.

    unify(N1, 2)
    
Destructuring. If Q1 holds a quant() relation, this equals, binds its arguments to new variables (`R1`).
    
    unify(Q1, quant(_, _, R1, _)

## not_equals

This function just compares two terms. If either their types or their values are unequal, it fails 

    not_equals(T1, T2)
    
* `T1`: a free variable, or any other term
* `T2`: a free variable, or any other term

## add

Adds two numbers `N1` and `N2` and places the result in `Sum`. If `Sum  is a number and it is not the sum of the arguments, the function fails.  

    add(N1, N2, Sum)
    
* `N1`: a number
* `N2`: a number
* `Sum`: a variable (to contain a number) or a number

## min

Sets `Min` to the smallest of `N1` and `N2.  

    min(N1, N2, Min)
    
* `N1`: a number
* `N2`: a number
* `Sum`: a variable (to contain a number) or a number

## compare

Adds two strings `N1` and `N2` and places the result in `R`. If N1 < N2, then R = -1; if N1 = N2, them R = 0; if N1 > N2, then R = 1;  

    compare(N1, N2, R)
    
* `N1`: a string
* `N2`: a string
* `R`: a variable (to contain a number)

This function is useful in order functions.

## subtract

Subtracts two numbers `N1` and `N2` and places the result in `Diff`. If `Diff  is a number and it is not the diff of the arguments, the function fails.  

    subtract(N1, N2, Diff)
    
* `N1`: a number
* `N2`: a number
* `Sum`: a variable (to contain a number) or a number

## date_today

Places the date of today in the variable in the form YYYY-mm-dd. If `D1` is not free, and does not contain today's date, the function fails.

* `D1`: a variable, to contain a date

    date_today(D1)

## date_subtract_years

Age calculation. If `D1` and `D2` contain dates, `Years` will be assigned the difference between these dates, in years (rounded down).  

    date_subtract_years(D1, D2, Years)
    
* `D1`: a variable, to contain a date
* `D2`: a variable, to contain a date
* `Years`: a variable, to contain a date
