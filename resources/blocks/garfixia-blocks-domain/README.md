# Blocks domain rules

The rules are separted into layers. Each layer has a specific responsibility. Each layer should only communicate with the layer directly below.

The layers are

- sentence
- command
- events  
- physics
- database

## Sentence layer

This layer picks up command as they are given in natural language. It performs the implicit actions that are expected by the user, but which are not named. These include getting rid of the object one is holding when one is told to pick up a new object. 

## Command layer

Once the command is cleaned up to be executed immediately, it is handled by the command layer. It performs the commands that are named explicitly in the sentence. This layer should be simple to understand, because it does exactly what you would expect it to do as a user.

## Events layer

All rules in this layer create events that can be used later for introspection ("when", "why", "how"). They pass around the id of the parent event.

## Physical layer

Commands have "physical" effects. Since this is just a simulation of a blocks world, all physical effects must be simulated. It is easy to forget about them in the command layer, and hard to debug. Therefore they have been given a separate layer. Anything that moves can change relationships between objects. This layer takes care of that. If you pick up a block, it may have been picked up from another block, so this removes the `support` relation and clears the top of the object below.

## Database layer

The database is only responsible for updating the database. It takes care of retracting and asserting knowledge.

## Cross layer

Each layer has a special rule file. Next to these files, there are two files that can be used by all layers:

- helper
- relations

