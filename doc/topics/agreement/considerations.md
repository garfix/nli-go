# Agreement - considerations

## Syntactic categories as sense?

The fact that categories are syntactic in nature means that we can't use them as a sense. 

So this won't work

    { rule: pronoun(E1) -> 'he',    sense: dom:male(E1) }
    { rule: pronoun(E1) -> 'it',    sense: dom:object(E1) }

In stead, we use this

    { rule: pronoun(E1) -> 'he',    tag: go:category(E1, gender, masculine) }
    { rule: pronoun(E1) -> 'it',    tag: go:category(E1, gender, neuter) }

In practise you would use both, the tag for agreement, the sense to restrict the query.
    
## Agreement within and between entities

To have subject and verb agree, we write

    { rule: vp(P1) -> np(E1) tv(P1, E1, E2) np(E2),      tag: go:agree(P1, E1) }

This syntax is not suitable to describe that a noun should agree with its modifiers, because they share the same variable. We could write `tag: go:agree(E1, E1)`, but instead we'll assume that the tagged categories of an entity should always match. So there's an implicit agreement of all categories within an entity.

