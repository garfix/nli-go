# Config files

The command line application needs a JSON config file to work. This config file specifies the resources used.
 
Check the _resources_ directory for example config.json files.

A relative path in a config files has that config file as its base.
 
## Sections

### grammars
 
 "grammars" is an array of file paths of grammars.
 A grammar is a set of grammar rules: rewrite rules, with their relational meaning representation (sense).
 
### generationgrammars
 
 "generationgrammars" is an array of file paths of grammars meant to create (generate) responses.
 
### factbases
 
 "factbases" is an array of databases of the following types

Each fact base may have a `sharedIds` field that holds the path to a JSON file. This file holds a mapping of local
database ids to shared database ids.
  
#### relation
 
 An array of "relation" fact bases.
 
 A "relation" fact base is read-only fact base specified by a .relation file that contains facts in a relational form.

 The fact base also has a .map file that specifies how a single domain specific relation matches to one or more relations in the fact base.

 Example:
 
     "relation": [
       {
         "name": "soil-samples",
         "facts": "relationships.relation",
         "map": "relationships-db.map"
       }
     ],

#### mysql
 
 An array of "mysql" fact bases.
  
 A "mysql" fact base is a preexisting read-write MySql database.
 
 This database may be used or unused (enabled) in this request, its access specifics are given (domain, username, password, database).
 The tables and columns used by the NLI app are named.
 
 The config entry also has a .map file that specifies how a single domain specific relation matches to one or more relations in the database (map).
 
 Example:
 
     "mysql": [
       {
         "enabled": false,
         "domain": "localhost",
         "username": "root",
         "password": "",
         "database": "my_nligo",
         "tables": [
           {
             "name": "marriages",
             "columns": [ { "name": "person1_id" }, { "name": "person2_id" }, { "name": "year" } ]
           },
           {
             "name": "parent",
             "columns": [ { "name": "parent_id" }, { "name": "child_id" } ]
           },
           {
             "name": "person",
             "columns": [ { "name": "person_id" }, { "name": "name" } ]
           }
         ],
         "map": "relationships-db.map"
       }
     ]
   }

#### sparql

    "sparql": [
      {
        "name": "dbpedia",
        "baseurl": "https://dbpedia.org/sparql",
        "defaultgraphuri": "http://dbpedia.org",
        "map": "dbpedia/ds2db.map",
        "names": "dbpedia/names.json",
        "entities": "dbpedia/entities.json",
      }

##### names

Sparql has the property 'names', which is for example:

    {
      "spouse": "http://dbpedia.org/ontology/spouse",
      "child": "http://dbpedia.org/ontology/child",
      "children": "http://dbpedia.org/property/children",
      "type": "http://www.w3.org/1999/02/22-rdf-syntax-ns#type",
      "description": "http://purl.org/dc/terms/description"
    }

This structure maps database names to URI's. This way we can talk to Sparql just like a relational database, with relations like 'spouse' ans 'description'. But when the actual Sparql query is created, these relations are turned into URI's.


Each of the fact bases can have the properties 'map', 'names', 'entities':

#### map

Map maps domain specific relations to database relations.

On the left side of the arrow you find the domain specific relations, and on the right side the domain specific relations. Example:

    [
        married_to(A, B) => spouse(A, B);
        married_to(A, B) => spouse(B, A);

        name(A, F, first_name) name(A, L, last_name) join(N, ' ', F, L) => birth_name(A, N);

        description(A, D) => description(A, D);

        gender(A, male) => gender(A, 'male');
        gender(A, female) => gender(A, 'female');

        person(E) => type(E, `http://dbpedia.org/ontology/Person`);
    ]

This map is used to create the database relations. It is also used to determine 'relation groups': groups of relations that need to stay together when used with a database. The left hand side of a mapping forms such a relation group.

#### entities

Entities are used to resolve proper names into database specific identifiers. Here's an example:

    {
      "person": {
        "name": "[name(Id, Name, full_name) person(Id)]",
        "knownby": {
          "description": "[description(Id, Value)]"
        }
      }
    }

Here "person" is a random identifier that will be shown to the end user if he/she needs to disambiguate the name. Next to "person" you can specify other entities.

"name" specifies a domain specific relation set. When used, the system will fill variable Name, query the relation set in the factbase. The bound Id variable will be stored as the id of the named entity.

"knownby" specifies some means to help a human user to disambiguate an entity. Here "description" is, like "person" a random identifier to be shown to the user.
Id is the variable that is entered by the system, just before the relation set is queried in the fact base. The Value variable is filled with data from the fact base.


In this example the fact base contained a relation "description(Id, Value)" that holds a proper description of a person. In other cases, multiple entity attributes may be needed to specify a person / entity to a user.

### rulebases
 
An array of paths to rule base specifications.
  
A rule base is a set of Datalog-like inference rules are be used by the system to infer new facts from existing facts, in order to solve a query.
 
### solutions
 
An array of paths to solutions.
 
A solution is the way the system uses to link a problem to a solution.
 
### predicates

The path of a predicates json file, which looks like this:

    {
        "has_capital": {"entityTypes": ["country", "city"] }
    }

This file contains domain specific predicates, the ones that are used in transformation files.

Here you can specify the entity types of the arguments.

These entity types are used for name resolution. If a name is used in the sentence, the system uses the entities file to look up the names. At the same time it will look at the relations and the predicates file.
From this it will find out what entity type belongs to the name. It will then only look for names that belong to this entity type.

It is optional to specify the predicates in this file. If there is no need to specify the entity types, they may be omitted.
