package mentalese

type FunctionBase interface {

	// Performs argument.Predicate on argument.Arguments[1..] and returns the new value of the first argument (0)
	Execute(input Relation, binding Binding) (Binding, bool)
}
