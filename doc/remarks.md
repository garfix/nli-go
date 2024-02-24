## 2024-02-24

The fact that a function-call term is implemented as a 1-relation relation set can get me into trouble. A relation set with one function-call that doesn't return any values is not treated as a function call.

I used to use the internal database to store 2d grids. I turned the 2d grid in a 1d list, and expected it to be faster, but in stead it got slower. This means that the overhead of each step is way too big at the moment and simple structures like [A * B] take up too much time and must be optmized.

## 2024-02-17

I created a function-specific version of if-then / if-then-else, called if2. I'm using this only in functions. Once all procedures with `if` have been rewritten to functions, if2 can be renamed to if, and the old if can go.

## 2024-02-11

Implemented returning zero values or multiple values. More work than I thought, but it's done. Yes!

If a function returns 1 value, it returns this value. If it returns multiple values, or no values, it returns a list. The number of return values is stored in a special field.

## 2024-02-07

Returning multiple values.

## 2024-02-04

Removing the [] brackets around assignments. In order to do this relations need to be visually grouped. So I'm adding {} around relations when they form a term. But since this happens to much, I make this optional. If that works it would be great.

## 2024-01-28

A simple solution would be to restore the previous state of the loop-variables after the loop.

Make a difference between immutable statements and mutable statements (difference variable treatments)?

What statements would then be mutable?

* for loops
* if statements
* assert / retract
* break / cancel / return

What statements could be both mutable and immutable?

* all operators
* assignment (needed for function calls in immutable environment)
* append
* log

I'm creating a scope by restoring the variables after the execution. The scope is not part of a data structure.

## 2024-01-21

Start a lexical scope within a for-loop: variables that are set here should loose there value when the scope is closed.

## 2024-01-20

Current problem: in a function all relations are mutable. When they are passed to `ExecuteChildStackFrame`, this will bring only 1 result, formed by the last values of the mutable variables.

## 2024-01-13

I intend to replace all custom "mentalese" programming code with plain Go. But this is a massive undertaking. I considered a sort of hybrid, but this would make things very complicated. Still, I need to do a lot of thinking before I can start.

How to execute something like this, in Go:

    go:do($np1, go:do($np2, dom:do_put_on_smart(E1, E2)))

Another option: use only local variables and just remove all non-local variables.

A reason why Go would not be a good fit is that most structures passed are declarative.

