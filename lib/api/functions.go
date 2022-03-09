package api

import "nli-go/lib/mentalese"

type SimpleFunction func(messenger SimpleMessenger, relation mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool)

type SolverFunction func(messenger ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet

type MultiBindingFunction func(messenger ProcessMessenger, goal mentalese.Relation, bindings mentalese.BindingSet) mentalese.BindingSet

type RuleFunction func(goal mentalese.Relation, binding mentalese.Binding) ([]mentalese.RelationSet, mentalese.BindingSet)
