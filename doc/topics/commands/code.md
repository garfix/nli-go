# Code

As commands have senses that diverge from the standard (declarative) interpretation, they have a completely separate set of rules.

Here are the basic rules for imperative sentences:

{ rule: s(S1) -> imperative(S1) }
{ rule: imperative(P1) -> imperative_clause(P1) '.' }
{ rule: imperative(P1) -> imperative_clause(P1) }
{ rule: imperative_clause(P1) -> vp_imperative(P1),         intent: go:intent(command) }

The `vp_imperative` has multiple instances, like:

    { rule: vp_imperative(P1) -> 'put' np(E1) 'in' np(E2),                            sense: go:do($np1, go:do($np2, dom:do_put_in_smart(E1, E2))) }

Note that a predicate responsible for executing the task is prefixed with `do_`, in order to distinguish it from declarative predicates, and to make it clear that its intent is to make a change to the world.

## go:do()

In the sense of a command we don't use `go:check`, we use `go:do`. The difference between them is that `go:check` applies its body to all members in the range of its quant, then checks if the number of members that succeed satisfies the quantifier. Whereas `go:do` applies its body to just enough members to satisfy the quantifier.

For example: we do `go:check()` on Fruit. Fruit has 5 instances: apple, banana, citrus, dade, and eggplant. The question we're asking is

    Are 3 fruits rotten?

This looks something like this

    go:check(go:quant(3, E, fruit(E)), is_rotten(E))

And let's suppose that only the apple, the banana, and the eggplant are rotten.

`go:check()` will then try to apply `is_rotten` to all 5 fruit. Apple, banana and eggplant (3 elements) survive this test. Then `go:check()` tests if 3 equals 3. This succeeds.

An example with `go:do()` would go something like:

    Throw away 1 rotten fruit.

This would look like something like this

    go:do(go:quant(1, E, fruit(E) is_rotten(E)), throw_away(E))

Now suppose that `go:do` would act just as `go:check`. It would consider each `fruit` that `is_rotten` and apply `throw_away` to it. This happens to 3 fruits. Finally it would check if 3 equals 1. This fails. So the predicate both fails, and 3 fruit are now in the trash can. This is obviously wrong.

What `go:do` actually does is this: it takes the `is_rotten` `fruit` apple, throws it away, increases the count to 1, and checks if the quantifier matches. 1 = 1, so `go:do` stops and succeeds.

`go:do` can't act on too many elements. It always succeeds, unless the number of elements it can act on does not reach the quantifier. In the example `go:do` would fail if there are no rotten fruits around.
