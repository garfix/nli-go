# Ontology

In your application you model a part of the world in your own, application specific way. It is up to you how you model
it; you do not have to worry if your model is somehow "right"; just make sure the application works well for your users
and you can extend it easily.

The semantic modelling of the world in an NLI system is called an ontology. It has several parts.

## Sorts

Any entity you are using has a type: the sort. Each sort needs a simple identifier, like

    person
    work
    country

The sorts themselves are not defined anywhere, but you will need them in several places.

## Relations

The relations of the ontology look simply like this

    person(Id, Name)
    write(PersonId, WorkId)

Although there are also relations in the databases you use, try to design the relations in the ontology independently.
The relations should model the way users think about the domain.

You also need to specify sort-properties.yml

    person:
      name: person_name(Id, Name)
      knownby:
        description: description(Id, Value)
    country:
      name: country_name(Id, Name)
      knownby:
        label: label(Id, Value)
        founding_date: founding_date(Id, Value)

You can see where the sorts pop up.

The "name" specifies the relations you need to look up a name (proper noun) in the database.

The "knownby" tells the system where to find a description of the entity. Used for disambiguation by the user.

## Id

An id identifies an entity and looks like this

    `person:123`

In general:

    `sort:shared-id`

When the entity only exists in a single database, `shared-id` is simply the primary key for the entity in the database.

If the entity occurs in multiple databases, with different id's, `shared-id` is an id from your application. You use
shared-id.json files to map from the application id to the database id. There's one such file for each database.

## Inheritance

In the future you would like to say: employer extends person.

Not implemented yet.