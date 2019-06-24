# A binding-processing function

The library contains many functions that have a relation set as input, together with a set of bindings and possibly a key cabinet.

The relation set is part of the query that needs answering and many be some levels deep in the processing flow.

For example:

aaa(X) bbb(Y) ccc(Z)

It may already have some bindings, but there are possibly none, and there may be variables in the binding that do not exist in the relation set.

For example

X=1 R=3

The key cabinet contains variable bindings that are knowledge base specific. They should only be used when the query is transformed into a knowledge base query.

dbpedia: { Y=http://dbpedia.nl/byron } in-memory: { S=823 }

At the same time, the function at hand may use a datastructure that itself is defined by a number of variables. These variables are from a different name space than the ones from the relation set, and should not be confused.

solution:
    pre: aaa(E) :- ddd(E, F) eee(F, G)
    post: fff(F)

Note that the variable F occurs both in the "pre" and "post" structure of this "solution".

Now, when this solution is used in the function, the whole structure must be converted to fit the input variables.

In this case E must be turned in X.

The other variable of the structure, F and G, must be converted into variables that are different from the input variables.

The result is something like this:

solution:
    pre: aaa(X) :- ddd(X, Q1) eee(Q1, Q2)
    post: fff(Q1)

Now the solution can be processed in terms of the original bindings and the key cabinet.

The resulting bindings are also from the namespace of the input variables, but they may contain some extra variables

X=1 R=3 Q1=18 Q2=23

The remaining variables must not leave the function because they have no function there.

Therefore the resulting bindings must be filtered with

- the variables from the relation set
- the variables from the input binding

## Skeletal function

Here's an example of what this looks like

    func (solver ProblemSolver) SolveSingleRelationSingleBindingSingleRuleBase(goalRelation mentalese.Relation, keyCabinet *ResolvedKeyCabinet, binding mentalese.Binding) mentalese.Bindings {

        inputVariables := goalRelation.GetVariableNames()

        goalBindings := mentalese.Bindings{}

        // find helper structure (named source)
        sourceSubgoalSets, sourceBindings := findSource(goalRelation)

        for i, sourceSubgoalSet := range sourceSubgoalSets {

            sourceBinding := sourceBindings[i]

            // rewrite the variables from the source namespace to our namespace
            importedSubgoalSet := sourceSubgoalSet.ImportBinding(sourceBinding)

            // perform the actual action
            subgoalResultBindings := solver.SolveRelationSet(importedSubgoalSet, keyCabinet, mentalese.Bindings{binding})

            // process the resulting bindings
            for _, subgoalResultBinding := range subgoalResultBindings {

                // filter out the input variables
                filteredBinding := subgoalResultBinding.FilterVariablesByName(inputVariables)

                // make sure all variables of the original binding are present
                goalBinding := binding.Merge(filteredBinding)

                goalBindings = append(goalBindings, goalBinding)
            }
        }

        solver.log.EndDebug("SolveSingleRelationSingleBindingSingleRuleBase", goalBindings)

        return goalBindings
    }
