# Sorts

The sort of an entity is its type or category. The sort of Mary is `person`. The sort of a table is `table`.

Sorts are declared in several files:

** argument-sort.relation ** lists the sorts of the arguments of each predicate.

    contains(container, object)
    loves(person, entity)

This is the basic tool that NLI-GO has to determine the sort of an entity. From the relations in the sentence, the system looks up the sorts in this file.

These are needed for name resolution and sortal filtering. 

To find a name in a database it is useful to know its sort. It reduces the search space.

** sort-properties.yml ** lists some properties of a sort.

    country:
        name: country_name(Id, Name)
        knownby:
            label: label(Id, Value)
            founding_date: founding_date(Id, Value)
        entity: country(Id)

These are needed to disambiguate names, and to solve one-anaphora.

To look up the id of an entity whose sort is `country`, this file is used. The relation `country_name()` can be used to find the name. If multiple ids are found, `knownby` is used to present the user descriptions of possible alternatives.

The `entity` is used for one-anaphora. Once the sort is known, the system needs to know which relation belongs to it. This relation is added to the relational representation of the sentence.

And finally

** sort-hierarchy.txt ** defines sort - subsort relations

    entity > area
    area > country
    area > state
    entity > city
    entity > person

The basic sort is `entity`. It is built-in. All other sorts derive from it.

The hierarchy is used to prevent conflicts that are not real.


