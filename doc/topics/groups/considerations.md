# Considerations

Each entity is represented by a variable in the dialog. The variable may hold multiple values.

## Multiple bindings or single binding with list?

A variable may have multiple bindings, and this is the best way to deal with groups while processing a sentence.

After the intent is executed, and an answer needs to be generated, it's easier to deal with a group as a list of ids. This way we only need to generate a sentence based on a single binding.

## Problem

When the sentence is complete, the group is added to the dialog bindings as a list. The next sentence then needs to deal with this list, but it doesn't. This is currently an unresolved issue.

