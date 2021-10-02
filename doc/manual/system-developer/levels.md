# Levels

For every language feature you need to find out on which level it belongs. 

## Code levels

These are the available levels:

* Built-in relations and structures (dialog context)
* System programmable rules (respond.rule) 
* Custom programmable rules

A built-in feature is automatically available to the user/programmer. You want as many features built-in so that the system supports the user as much as possible. However, features that are not universal should not be built-in. Building them in causes confusion when they are not wanted.

## Data levels

If a feature stores data, it must be decided on which level this is done:

* Database (built-in or custom)
* Dialog Context (built-in)
* Process (built-in)

Storing data in the database is permanent. Data stored in the dialog context is available to the rest of the dialog. Data stored in the process is accessible to the rest of the process.  

## The levels of a feature

A language feature is ideally implemented on one code level and one data level. This is conceptually easiest.

However, features such as anaphora are stored by code on the built-in relation level and the system programmable levels.
