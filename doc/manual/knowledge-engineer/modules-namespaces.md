# Modules and namespaces

An application is rooted in its `config.yml` file that looks like this:

    uses:
      dom: garfixia-dbpedia-domain:1.0.0
      dbpedia: garfixia-dbpedia-db:1.0.0
      en: garfixia-dbpedia-en-us:1.0.0
      sol: garfixia-dbpedia-solution:1.0.0

The config describes the modules of the application. Each line has a module. In

    dom: garfixia-dbpedia-domain:1.0.0
    
`dom` is an alias, `garfixia-dbpedia-domain` is the name of the module (and its directory name), and `1.0.0` is its required version.

The alias is used by the system to prefix relations, in order to keep same relations from different modules apart. It is also used to name database connections; in the case of `dbpedia` for instance. Use short but meaningful aliases.

These aliases are called `application aliases`, to distinguish them from `module aliases` that we will see later.

Use of versions is currently very simple: the version either matches or it doesn't. I will extend this when the need arises.

## Module types

There are 4 types of modules:
    
    * domain
    * solution
    * grammar
    * db
    
Domain modules contain rules, predicates, and sort hierarchy.

Solution modules just contain solution files.

Grammar modules contain read and write rule sets.

Database modules contain specifications of MySQL, SPARQL and Internal databases.

## Use of namespaces

Each module has an `index.yml` file that describes the external modules it uses:

    uses:
      dom: garfixia-dbpedia-domain:1.0.0

Here `dom` is a module alias that is used within this module only. They are used as prefixes. When a relation of this external module is used, this looks like

    dom:has_population(E1, E2)
    
Relations that are defined by the module itself are not prefixed.    

These aliases are called `module aliases`. They are used by NLI-GO to locate relations in specific modules, but you will not see them when you execute a request. They are only used to build the system. When the system executes a request, it's only the application aliases that you will see in a debug trace.  

## The go alias

Built-in predicates are always prefixed by the system-alias `go`. For example:

    go:check($np, dom:has_father(E2, E1))
    
    