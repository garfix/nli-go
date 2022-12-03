# Sorts

The type of an entity is called its "sort", traditionally. The sort of Mary is `person`, for example. The sort of a table is `table`.

## Declare the sorts of database entities

You can define the sort of an entity by declaring it explicitly:

    go:has_sort(`block:big-red`, block)

Commonly, the sort of an object is stored in a database, and the map-file declares how they are mapped. 
In the blocks world, for example, we get:

    go:has_sort(E, T) :- is(E, T);

`is()` is a relation in the database that contains sort information.

## Sort hierarchy

A sort can have a sort of itself. Here's how we declare a block to be an object:

    go:has_sort(block, object)

We can declare a snake to be a reptile, and both a reptile and a mammal to be animals: 

    go:has_sort(snake, reptile);
    go:has_sort(reptile, animal);
    go:has_sort(mammal, animal);

Thus we can create a "sort hierarchy". Such a hierarchy has a single root sort (`entity`), and any number of subsorts.

    go:has_sort(animal, entity);

While it is common to create a simple hierarchy with sorts, it is not strictly necessary. It's possible for an entity or a sort to have multiple supersorts:

    go:has_sort(`cartoon:tweety`, bird);
    go:has_sort(`cartoon:tweety`, cartoon_character);

## Declare the sorts of sentence entities

NLI-GO needs to know the sorts of some of the entities before it can process the sentence:

- if the sentence contains a name, the system needs to know what sort it has, in order to efficiently look up the name in the database
- anaphora resolution needs the type of entities as well: to avoid sort conflict, and to make sure the antecedent of a reference has a sort
- sortal filtering weeds out sentences that contain an entity with conflicting sorts [1]

[1] Note that this can only work with a sort hierarchy, in which an entity has only a single sort. Sort conflicts can't be detected if entities are allowed to have multiple types 

NLI-GO makes use of the `go:has_sort()` senses in the grammar, to induce the sort of the entity.

Another source of information is the file `argument-sort.relation` It lists the sorts of the arguments of each predicate:

    contains(container, object)
    loves(person, entity)

This declaration of argument sort is part of the domain definition. They are not defined with each word in the grammar, because the grammar may contain the relation more than once. Also, the sort is language-independent. It does not need to be declared for each separate language.

To find a name in a database it is useful to know its sort. It reduces the search space. The file `sort-properties.yml` lists some properties of a sort.

    country:
        name: country_name(Id, Name)
        knownby:
            label: label(Id, Value)
            founding_date: founding_date(Id, Value)
        entity: country(Id)

These are needed to disambiguate names, and to solve one-anaphora.

To look up the id of an entity whose sort is `country`, this file is used. The relation `country_name()` can be used to find the name. If multiple ids are found, `knownby` is used to present the user descriptions of possible alternatives.

The `entity` is used for one-anaphora. Once the sort is known, the system needs to know which relation belongs to it. This relation is added to the relational representation of the sentence.

## go:isa()

The predicate `go:isa()` can be used to determine if entity E has sort S, either directly, or indirectly.

## Each id has a sort

Each id, commonly an integer, is a unique identifier, but only within its own sort. Therefore the sort is added to the id, on creation.

Hence it is also possible to determine the sort from a given id. 