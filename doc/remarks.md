## 2024-01-20

Current problem: in a function all relations are mutable. When they are passed to `ExecuteChildStackFrame`, this will bring only 1 result, formed by the last values of the mutable variables.

## 2024-01-13

I intend to replace all custom "mentalese" programming code with plain Go. But this is a massive undertaking. I considered a sort of hybrid, but this would make things very complicated. Still, I need to do a lot of thinking before I can start.

How to execute something like this, in Go:

    go:do($np1, go:do($np2, dom:do_put_on_smart(E1, E2)))

Another option: use only local variables and just remove all non-local variables.

A reason why Go would not be a good fit is that most structures passed are declarative.

