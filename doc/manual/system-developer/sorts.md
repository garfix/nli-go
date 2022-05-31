# Sorts

The sort of an entity is its type or category. The sort of Mary is `person`. The sort of a table is `table`.

Sorts are declared in several files:

** argument-sort.relation ** lists the sorts of the arguments of each predicate.

    contains(container, object)
    loves(person, entity)

These are needed for sortal filtering.

** sort-properties.yml ** lists some properties of a sort.

    country:
        name: country_name(Id, Name)
        knownby:
            label: label(Id, Value)
            founding_date: founding_date(Id, Value)
        entity: country(Id)

These are needed to disambiguate names, and to solve one-anaphora.

And finally

** sort-hierarchy.txt ** defines sort - subsort relations

    entity > area
    area > country
    area > state
    entity > city
    entity > person

The basic sort is `entity`. It is built-in. All other sorts derive from it.

The hierarchy is used to prevent conflicts that are not real.


