# Config files

The command line application needs a JSON config file to work. This config file specifies the resources used.
 
Check the _resources_ directory for example config.json files.

A relative path in a config files has that config file as its base.
 
## Sections
 
### lexicons
 
 "lexicons" is an array of file paths of lexicons. 
 A lexicon is a set of lexical items: words, with their parts-of-speech (pos) and relational meaning representation (sense).
 
### grammars
 
 "grammars" is an array of file paths of grammars.
 A grammar is a set of grammar rules: rewrite rules, with their relational meaning representation (sense).
 
### generationlexicons
 
 "generationlexicons" is an array of file paths of lexicons meant to create (generate) responses.
 
### generationgrammars
 
 "generationgrammars" is an array of file paths of grammars meant to create (generate) responses.
 
### factbases
 
 "factbases" is an array of databases of the following types
  
#### relation
 
 An array of "relation" fact bases.
 
 A "relation" fact base is read-only fact base specified by a .relation file that contains facts in a relational form.

 The fact base also has a .map file that specifies how a single domain specific relation matches to one or more relations in the fact base.

 Example:
 
     "relation": [
       {
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
        "baseurl": "https://dbpedia.org/sparql",
        "defaultgraphuri": "http://dbpedia.org",
        "map": "dbpedia-db.map",
        "names": "dbpedia-db-names.json",
        "stats": "dbpedia-db-stats.json"
      }

### rulebases
 
 An array of paths to rule base specifications.
  
 A rule base is a set of Prolog-like inference rules are be used by the system to infer new facts from existing facts, in order to solve a query. 
 
### solutions
 
 An array of paths to solutions.
 
 A solution is the way the system uses to link a problem to a solution.  
 
### generic2ds
 
 An array of paths to generic-to-domain-specific transformations.
 
### ds2generic
 
 An array of paths to domain-specific-to-generic transformations.
 