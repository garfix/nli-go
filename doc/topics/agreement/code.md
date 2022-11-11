# Agreement - code

It's important to note that agreement annotations are not necessary to interpret most sentences correctly. If the input contains "the children picks up a ball", this is syntactically incorrect, but may still be transformed into a proper sense. Try to build a grammar without agreement annotations first, and only when you notice that it leads to undesired effects, add them.

When no parse can be made, due to an agreement conflict, the system will reply with a statement like:

    "Agreement mismatch: plural / singular"

`gender` and `number` have a built-in meaning in the system. Other syntactic categories can be created at will.

## Person

Make the subject and the verb agree on person

"A boy steals some magazines"

{ rule: vp(P1) -> np(E1) tv(P1, E1, E2) np(E2),                                 tag: go:agree(P1, E1) }
{ rule: tv(P1, E1, E2) -> 'steals',                                             tag: go:category(P1, person, third) }
{ rule: noun(E1) -> 'boy',                                                      tag: go:category(E1, person, third) }

"Some magazines were stolen by a boy"

{ rule: vp(P1) -> np(E1) aux_be(P1) past_participle(P1, E2, E1) 'by' np(E2),    tag: go:agree(P1, E1) }
{ rule: aux_be(P1) -> 'were',                                                   tag: go:category(P1, person, third) }
{ rule: noun(E1) -> 'boy',                                                      tag: go:category(E1, person, third) }

## Number

Make the subject and the verb agree on number

"A boy steals some magazines"

{ rule: vp(P1) -> np(E1) tv(P1, E1, E2) np(E2),                                 tag: go:agree(P1, E1) }
{ rule: tv(P1, E1, E2) -> 'steals',                                             tag: go:category(P1, number, singular) }
{ rule: noun(E1) -> 'boy',                                                      tag: go:category(E1, number, singular) }

## Categories based on names

Proper names are looked up in the database. When a name is found, the system can assign syntactic categories to the entity with the name. In order to do this, use the file `sort-properties.yml`:

    person:
        name: name(Id, Name)
        gender: gender(Id, Value)
        number: [Value := 1]

Read this as follows: for the sort `person`, if a name is detected using the relation `name`, the value from Id is bound to the entity variable. When that happens, also try the relation `gender`; if it succeeds, `Value` will contain the gender and this will be bound as `gender` category to the variable. In this example `number` is always bound to the value `1`, but this could be replaced with another relation.
