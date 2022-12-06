# Code

## Noun

A noun is represented by the sort it embodies:

{ rule: noun(E1) -> 'pyramid',                                         sense: go:has_sort(E1, pyramid) }
{ rule: noun(E1) -> 'pyramids',                                        sense: go:has_sort(E1, pyramid) }
{ rule: noun(E1) -> 'block',                                           sense: go:has_sort(E1, block) }
{ rule: noun(E1) -> 'blocks',                                          sense: go:has_sort(E1, block) }
{ rule: noun(E1) -> 'cube',                                            sense: go:has_sort(E1, block) dom:cubed(E1) }
{ rule: noun(E1) -> 'cubes',                                           sense: go:has_sort(E1, block) dom:cubed(E1) }

`has_sort` is a built-in predicate that is used for different sortal restrictions.

## Adjective

An adjective, a modifier to a noun, has a simple sense:

{ rule: adjective(E1) -> 'blue',                                       sense: dom:blue(E1) }
{ rule: adjective(E1) -> 'green',                                      sense: dom:green(E1) }
{ rule: adjective(E1) -> 'big',                                        sense: dom:big(E1) }
{ rule: adjective(E1) -> 'small',                                      sense: dom:small(E1) }

## Verb 

A verb is represented by a single predicate. There are different categories.

Simple transitive verb (subject, object)

    { rule: tv(P1, E1, E2) -> 'support',                                   sense: dom:support(E1, E2) }

`P1` represents the event in which the predicate is involved. More general, it is a predication, that predicates ("says") something about one or more entities. "Predication" is abbreviated to "P" in the variable name.

If the application suggests that is easier to split up the predicate into more specific predicates, this is also possible:

    { rule: tv(P1, E1, E2) -> 'support',                               sense: dom:support(P1) dom:supporter(P1, E1) dom:supportee(P1, E2)}

However, this approach involves more work to the developer, and is only justified if the application benefits from it.

Non-finite have a specific representation. Here we have *infinitive* and *gerund*

    { rule: tv_infinitive(P1, E1, E2) -> 'pick' 'up',                  sense: dom:pick_up(P1, E1, E2) }
    { rule: tv_gerund(P1, E1, E2) -> 'holding',                        sense: dom:hold(P1, E1, E2) }

The "pick up" example shows that a verb may have a particle attached ("up").

These forms occur in specific places in other rules.

## Adverb

The adverb, a modifier to a verb, is a modifier to a predication. It can be used to provide an attribute to a command, for example.

    { rule: adjective(P1) -> 'slowly',                                 sense: dom:slowly(P1) }

## Preposition

Prepositions like "in" and "on top of", are represented by relations. There's no need for a `preposition` category, the `pp` (preposition phrase) does the job.

    { rule: pp(E1) -> 'in' np(E2),                                     sense: go:check($np, dom:contain(_, E2, E1)) }
    { rule: pp(E1) -> 'on' 'top' 'of' np(E2),                          sense: go:check($np, dom:on(E1, E2)) }

But when they are part of a command, they can best be represented like they belong to the verb:

    { rule: vp_imperative(P1) -> 'put' np(E1) 'on' np(E2),             sense: go:do($np1, go:do($np2, dom:do_put_on_smart(E1, E2))) }

The command "put in" simply requires a different action than "put on". There's no use for the abstract "put".

## Conjunction

The sense of the conjunction "and" coincides with the default operation that combines two senses. That is, there is no need to explicitly *and* two senses, because they would be *anded* without the operator as well. So NLI-GO has no sense for "and", it's implicit.

Having said that, it can be useful to represent "and" with the sense `go:and()` in some cases.

    { rule: np(E1) -> np(E1) 'and' np(E2),                                 sense: go:and($np1, $np2) }

In the blocks world, this is needed because a compound nominal, consisting of "and" and "or"'s, is passed to the predicate `go:quant_ordered_list()`.