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

