# Linguistics

In the response to a question, the system usually needs to name the entities it found in the search.

## Persons

If an entity is a person, the name of the person can be used, of course. If person can have multiple names, a canonical name should be available for the response. For example:

Available names:

- Lord Byron
- George Gordon Byron

Canonical name: "Lord Byron".

## Objects

In a set of objects on a table, the objects don't have individual names. The objects need to be described from their attributes and relations. The attributes are predefined. Use as little attributes as necessary. There's a very specific order in which to add attributes,

Every object has a form, and the form should be as specific as possible. Choose "cube" over "block". If there's only one object in the scene with this form, this is enough. For example: box. This uniquely determines the object, and thus we say "the box".

If there are multiple objects with this form, add an attribute: color. If the combination of shape and form is unique, stop. We can say "the green pyramid".

If there are multiple objects with this form and color, add an attribute: volumne. If the combination of shape, form and volume is unique, stop. We can say "the large red block".

If there are multiple objects with this form, color and volume, add a relation: support. This relation links to another object. SHRDLU just presumes that this combination of shape, form, volume and support is unique. "the large green block which supports the red pyramid"

todo: why is this answer using indefinite articles?
{"Is there a large block behind a pyramid?", "Yes, three of them: a large red one, a large green cube and the blue one"},

## Grouping

When two or more object descriptions are the same, they are grouped together. In stead of "the large green cube and the large green cube" we say "two large green cubes". Note the number "two" and the pluralization. Also, a new entity is introduced for the group; it replaces the instances.

## Definiteness

A question based on an indefinite NP will receive the specific entity in the response.

    Had you touched any pyramid before you put the green one on the little cube?
    Yes, the green one

whereas a question with a definite NP will just receive yes/no

    Have you picked up superblock since we began?
    Yes
