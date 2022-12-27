# Config files

An application has a single `config.yml` file in the root, and an `index.yml` file for each module.

All configuration is done in YAML, a low-syntax configuration language. If you're not familiar, you may want to [read about it](https://blog.stackpath.com/yaml/).

## config.yml

The config.yml file in the root of the application looks like this

    uses:
      dom: garfixia-dbpedia-domain:1.0.0
      dbpedia: garfixia-dbpedia-db:1.0.0
      en: garfixia-dbpedia-en-us:1.0.0
      sol: garfixia-dbpedia-solution:1.0.0

It has only a single entry, `uses`, which names the modules of the application.

In this example, `dom` is an alias, `garfixia-dbpedia-domain` is the name of the module (and its directory name), and `1.0.0` is its required version.

The alias is used by the system to prefix relations, in order to keep same relations from different modules apart. It is also used to name database connections; in the case of `dbpedia` for instance. Use short but meaningful aliases.

These aliases are called `application aliases`, to distinguish them from `module aliases` that we will see later.

Use of versions is currently very simple: the version either matches or it doesn't. I will extend this when the need arises.

## index.yml

There are 4 types of modules:
    
    * domain
    * solution
    * grammar
    * db
    
The `index.yml` file of a module starts with this:

    type: grammar
    version: 1.0.0
    uses:
      dom: garfixia-dbpedia-domain:1.0.0

Extra fields are added by different module types.    
    
Each type of module has a `type` field that names the kind of module, and a `version` that holds the version. Until version management becomes necessary, we use only 1.0.0.

All modules have a `uses` field through which relations of other modules are imported.
      
Here `dom` is a module alias that is only used within this module, to designate relations from external modules, by prefixing them with this alias.           
    
When a relation has its source in the module, it does not need to be prefixed by an alias. Prefix only external relations.     
    
### Domain modules

A domain module contains the knowledge about a single domain, like blocks, books, or some encyclopedia.

A sample index.yml:

    rules: [dbpedia.rule]
    write: [ write.yml ]
    sorts: sort-properties.yml
    
All of these are optional.     

`rules` contains the names of inference rule files. Inference rules are mostly rules like this:

    taller(A, B) :- height(A, Ha) height(B, Hb) go:greater_than(Ha, Hb);
    
but facts are also allowed:

    taller(`jack`, `john`);    

Currently sorts are only needed to locate proper names in databases. It helps to know that "Madonna" is a person, to avoid confusion with other entities with the same name. The sort is used in combination with the `entities` file of the database domain.

`write` has yml files that name the predicates of the goals of rules that can be written to this rule base. 

Here you can specify the sorts of the arguments. These sorts are used for name resolution. If a name is used in the sentence, the system uses sort-properties.yml to look up the names. At the same time it will look at the relations and the predicates file.
From this it will find out what sort belongs to the name. It will then only look for names that belong to this sort.

It is optional to specify the predicates in this file. If there is no need to specify the sorts, they may be omitted.

`entities` are used to resolve proper names into database specific identifiers. Here's an example:

    person:
        name: name(Id, Name, full_name) person(Id)
        knownby:
          description: description(Id, Value)

Here `person` is a sort.

`name` specifies a relation. When used, the system will fill variable Name, query the relation set in the factbase. The bound Id variable will be stored as the id of the named entity.

`knownby` specifies some means to help a human user to disambiguate an entity. Here "description" is, like "person" a random identifier to be shown to the user.
Id is the variable that is entered by the system, just before the relation is queried in the fact base. The Value variable is filled with data from the fact base.

In this example the fact base contained a relation "description(Id, Value)" that holds a proper description of a person. In other cases, multiple entity attributes may be needed to specify a person / entity to a user.

### Solution modules

A solution is the way the system uses to link a problem to a solution.

A sample index.yml:

    solution: [ dbpedia.solution ]
    
`solution` contains the names of [solution files](solution.md).    

### Grammar modules

A grammar module contains the read and write rules to parse and generate sentences.

A sample index.yml:

    read: [ read.grammar ]
    write: [ write.grammar ]
    tokenexpression: ([_0-9a-zA-Z]+|[^\\s])
    text: text.csv
    
`read` and `write` contain names of files that hold [grammar rules](creating-a-grammar.md).    

You can name the regular expression used to tokize a sentence if you are not satisfied with the standard expression.

`text` is an optional CSV file that maps source strings to translations.

These texts are used by `go:translate()`.

The CSV file contains lines like:

```
red,Rot
yellow,Gelb
blue,Blau
```

Quite simple. Some notes:

- all texts are trimmed of leading and trailing spaces and tabs
- quote are not allowed  
- if a text contains a comma, it must be escaped by a backslash, for example

```
red\, yellow\, and blue,Rot\, Geld und Blau
``` 

### Db modules    

Three types of database are supported; each has its own entries in the index.yml.

Common fields in index.yml:

    shared: [ shared-id.yml ]
    read: [ read.map ]
    write: [ write.map ]

Each fact base may have a `shared` field that holds the path to a Shared ids YAML file. This file holds a mapping of local
database ids to shared database ids.
    
    person:
      x21-r01: person1
      x21-r11: person3

Here `person` is a sort from the domain. `x21-r01` is an identifier from the database, and `person` is a shared identifier. When shared ids are used, NLI-GO reasons with these shared ids, and will convert them to database ids only when contacting the database. This allows you to combine the results of multiple databases in a single query.

In index.yml `read` and `write` are names of files that map from domain relations to database relations, when reading from or writing to the database.

Here are some examples from dbpedia's read file:

    dom:person_name(A, N) :- birth_name(A, N) type(A, `:http://dbpedia.org/ontology/Person`);
    dom:person_name(A, N) :- foaf_name(A, N) type(A, `:http://dbpedia.org/ontology/Person`); 
 
I will continue with the different types of database modules.

#### internal
 
 An internal fact base is a read/write fact base specified by a .relation file that contains facts in a relational form.

 Example index.yml:

    type: db/internal 
    facts: [ relationships.relation ]

`facts` holds the filenames of the files that hold relations that form the internal database.

#### mysql
 
 A MySql fact base is a connection to a read-write MySql database.
 
 This database may be used or unused (enabled) in this request, its access specifics are given (domain, username, password, database).
 The tables and columns used by the NLI app are named.
 
 The config entry also has a .map file that specifies how a single domain specific relation matches to one or more relations in the database (map).
 
 Example index.yml:
 
     type: db/mysql
     username: root
     password": root,
     database: my_nligo
     tables:
       -
         name: marriages
         columns: [ { name: person1_id }, { name: person2_id }, { name: year } ]
       -
         name: parent
         columns: [ { name: parent_id }, { name: child_id } ]
       -
         name: person
         columns: [ { name: person_id }, { name: name }, { name: gender }, { name: birthyear } ]

`database` is a database on the local machine. `username` and `password` allows us to log in.

`tables` holds the table names and columns of the tables that are relevant to the application. 

#### Sparql

A read-only SPARQL database, like DBpedia.

Example index.yml:

    type: db/sparql
    baseurl: https://dbpedia.org/sparql
    defaultgraphuri: http://dbpedia.org
    names: dbpedia/names.yml
    cache: true

Sparql has the property `names`, which is the name of a mapping file. This file maps database predicates to SPARQL urls, since in SPARQL each relation is identified by a URI.

    spouse: http://dbpedia.org/ontology/spouse
    child: http://dbpedia.org/ontology/child
    children: http://dbpedia.org/property/children
    type: http://www.w3.org/1999/02/22-rdf-syntax-ns#type
    description: http://purl.org/dc/terms/description

This structure maps database names to URI's. This way we can talk to Sparql just like a relational database, with relations like 'spouse' ans 'description'. But when the actual Sparql query is created, these relations are turned into URI's.

When `cache` is true, queries are cached and a single query is only made once. This speeds up repetitive queries tremendously. 

Currently the caching is forever, until the cache is cleared manually.  