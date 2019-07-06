package mentalese

import "nli-go/lib/common"

type SequenceApplier struct {
	log *common.SystemLog
}

// Turns the relation set
//
// abc(P1) sequence(P1, P2) def(P2)
//
// into the relation set
//
// seq([abc(P1)], [def(P1)])
func NewSequenceApplier(log *common.SystemLog) *SequenceApplier {
	return &SequenceApplier{log: log}
}

func (applier SequenceApplier) ApplySequences(set RelationSet) RelationSet {

	var sequenceRelations = RelationSet{}
	var remainingRelations = RelationSet{}

	for _, relation := range set {
		if relation.Predicate == PredicateSequence {
			sequenceRelations = append(sequenceRelations, relation)
		} else {
			remainingRelations = append(remainingRelations, relation)
		}
	}

	var seqRelations = RelationSet{}
	for _, relation := range sequenceRelations {
		var var0 = relation.Arguments[0]
		var var1 = relation.Arguments[1]
		var var2 = relation.Arguments[2]

		var0Relations := remainingRelations.findRelationsStartingWithVariable(var0.TermValue)
		remainingRelations = remainingRelations.RemoveRelations(var0Relations)
		var2Relations := remainingRelations.findRelationsStartingWithVariable(var2.TermValue)
		remainingRelations = remainingRelations.RemoveRelations(var2Relations)

		seqRelation := NewRelation(PredicateSeq, []Term{
			NewRelationSet(var0Relations),
			var1,
			NewRelationSet(var2Relations),
		})

		seqRelations = append(seqRelations, seqRelation)
	}

	newRemainingRelations := RelationSet{}
	for _, relation := range remainingRelations {
		if relation.Predicate == PredicateQuant {
			relation.Arguments[QuantScopeIndex] = NewRelationSet(applier.ApplySequences(relation.Arguments[QuantScopeIndex].TermValueRelationSet))
		}

		newRemainingRelations = append(newRemainingRelations, relation)
	}

	var resultSet = append(seqRelations, newRemainingRelations...)
	return resultSet
}