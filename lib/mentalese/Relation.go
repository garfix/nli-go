package mentalese

import (
	"nli-go/lib/common"
	"strings"
)

type Relation struct {
	Negate    bool   `json:"negate,omitempty"`
	Predicate string `json:"predicate"`
	Arguments []Term `json:"arguments"`
}

const Terminal = "terminal"

const ProcessInstructionBreak = "break"
const ProcessInstructionCancel = "cancel"
const ProcessInstructionReturn = "return"

const FrameTypePlain = "plain"
const FrameTypeScope = "scope"
const FrameTypeLoop = "loop"
const FrameTypeComplex = "complex"

const SortEntity = "entity"
const SortEvent = "event"

const PredicateCanned = "go_canned"
const PredicateQuantCheck = "go_quant_check"
const PredicateQuantForeach = "go_quant_foreach"
const PredicateQuantOrderedList = "go_quant_ordered_list"
const PredicateQuant = "go_quant"
const PredicateQuantifier = "go_quantifier"
const PredicateMakeList = "go_make_list"
const PredicateListAppend = "go_list_append"
const PredicateListOrder = "go_list_order"
const PredicateListForeach = "go_list_foreach"
const PredicateListDeduplicate = "go_list_deduplicate"
const PredicateListSort = "go_list_sort"
const PredicateListIndex = "go_list_index"
const PredicateListGet = "go_list_get"
const PredicateListHead = "go_list_head"
const PredicateListLength = "go_list_length"
const PredicateListExpand = "go_list_expand"
const PredicateRangeForeach = "go_range_foreach"
const PredicateAnd = "go_and"
const PredicateNot = "go_not"
const PredicateOr = "go_or"
const PredicateXor = "go_xor"
const PredicateCall = "go_call"
const PredicateIgnore = "go_ignore"
const PredicateAssert = "go_assert"
const PredicateRetract = "go_retract"
const PredicateIntent = "go_intent"
const PredicateEventReference = "go_event_reference"
const PredicateCount = "go_count"
const PredicateFirst = "go_first"
const PredicateLast = "go_last"
const PredicateGet = "go_get"
const PredicateOrder = "go_order"
const PredicateLargest = "go_largest"
const PredicateSmallest = "go_smallest"
const PredicateExists = "go_exists"
const PredicateCut = "go_cut"
const PredicateExec = "go_exec"
const PredicateExecResponse = "go_exec_response"
const PredicateSplit = "go_split"
const PredicateJoin = "go_join"
const PredicateConcat = "go_concat"
const PredicateCompare = "go_compare"
const PredicateUnify = "go_unify"
const PredicateMin = "go_min"
const PredicateDateToday = "go_date_today"
const PredicateDateSubtractYears = "go_date_subtract_years"
const PredicateSem = "go_sem"
const PredicateLog = "go_log"
const PredicateWaitFor = "go_wait_for"
const PredicateCreateGoal = "go_create_goal"
const PredicateUuid = "go_uuid"
const PredicatePrint = "go_print"
const PredicateFindLocale = "go_find_locale"
const PredicateSlot = "go_slot"
const PredicateIsa = "go_isa"
const PredicateGetSort = "go_get_sort"

const PredicateRespond = "go_respond"
const PredicateTokenize = "go_tokenize"
const PredicateParse = "go_parse"
const PredicateDialogize = "go_dialogize"
const PredicateEllipsize = "go_ellipsize"
const PredicateRelationize = "go_relationize"
const PredicateResolveNames = "go_resolve_names"
const PredicateCheckAgreement = "go_check_agreement"
const PredicateSortalFiltering = "go_sortal_filtering"
const PredicateResolveAnaphora = "go_resolve_anaphora"
const PredicateExtractRootClauses = "go_extract_root_clauses"
const PredicateDialogAddRootClause = "go_dialog_add_root_clause"
const PredicateDialogUpdateCenter = "go_dialog_update_center"
const PredicateDialogGetCenter = "go_dialog_get_center"
const PredicateDialogSetCenter = "go_dialog_set_center"
const PredicateGenerate = "go_generate"
const PredicateSurface = "go_surface"
const PredicateTranslate = "go_translate"

const PredicateDetectIntent = "go_detect_intent"
const PredicateSolve = "go_solve"
const PredicateFindResponse = "go_find_response"
const PredicateCreateAnswer = "go_create_answer"
const PredicateCreateCanned = "go_create_canned"

const PredicateUserSelect = "go_user_select"

const PredicateAlreadyGenerated = "go_already_generated"

const PredicateContextSet = "go_context_set"
const PredicateContextExtend = "go_context_extend"
const PredicateContextCall = "go_context_call"
const PredicateContextClear = "go_context_clear"

const PredicateDialogReadBindings = "go_dialog_read_bindings"
const PredicateDialogWriteBindings = "go_dialog_write_bindings"
const PredicateDialogAddResponseClause = "go_dialog_add_response_clause"

// internal relational representations of syntactic structures
const PredicateIncludeRelations = "$go$_include_relations"
const PredicateIfThen = "go_$if_then"
const PredicateIfThenElse = "go_$if_then_else"
const PredicateFail = "go_$fail"
const PredicateReturn = "go_$return"
const PredicateBreak = "go_$break"
const PredicateCancel = "go_$cancel"
const PredicateAssign = "go_$assign"
const PredicateEquals = "go_$equals"
const PredicateNotEquals = "go_$not_equals"
const PredicateGreaterThan = "go_$greater_than"
const PredicateLessThan = "go_$less_than"
const PredicateGreaterThanEquals = "go_$greater_than_equals"
const PredicateLessThanEquals = "go_$less_than_equals"
const PredicateAdd = "go_$add"
const PredicateSubtract = "go_$subtract"
const PredicateMultiply = "go_$multiply"
const PredicateDivide = "go_$divide"

const CategoryText = "text"
const CategoryProperNoun = "proper_noun"
const CategoryProperNounGroup = "proper_noun_group"

const TagRootClause = "go_root_clause"
const TagFunction = "go_function"
const TagAgree = "go_agree"
const TagReference = "go_reference"
const TagSortalReference = "go_sortal_reference"
const TagReflectiveReference = "go_reflective_reference"

const AtomFunctionSubject = "subject"
const AtomFunctionObject = "object"
const AtomFunctionNone = "none"

const AtomGender = "gender"

const AtomNone = "none"
const AtomSome = "some"
const AtomOne = "one"
const AtomAsc = "asc"
const AtomDesc = "desc"
const AtomReturnValue = "rv"

const QuantifierResultCountVariableIndex = 0
const QuantifierRangeCountVariableIndex = 1
const QuantifierSetIndex = 2

const QuantQuantifierIndex = 0
const QuantRangeVariableIndex = 1
const QuantRangeSetIndex = 2

const SeqFirstOperandIndex = 0
const SeqSecondOperandIndex = 1

const NotScopeIndex = 0

func NewRelation(negate bool, predicate string, arguments []Term) Relation {
	return Relation{
		Negate:    negate,
		Predicate: predicate,
		Arguments: arguments,
	}
}

func (relation Relation) GetPredicateWithoutNamespace() string {
	parts := strings.Split(relation.Predicate, "_")
	if len(parts) == 1 {
		return parts[0]
	} else {
		return relation.Predicate[len(parts[0])+1:]
	}
}

func (relation Relation) GetVariableNames() []string {

	var names []string

	for _, argument := range relation.Arguments {
		names = append(names, argument.GetVariableNames()...)
	}

	return common.StringArrayDeduplicate(names)
}

func (relation Relation) Equals(otherRelation Relation) bool {

	equals := relation.Predicate == otherRelation.Predicate

	equals = equals && relation.Negate == otherRelation.Negate

	for i, argument := range relation.Arguments {
		equals = equals && argument.Equals(otherRelation.Arguments[i])
	}

	return equals
}

func (relation Relation) Copy() Relation {

	newRelation := Relation{}
	newRelation.Predicate = relation.Predicate
	newRelation.Negate = relation.Negate
	newRelation.Arguments = []Term{}
	for _, argument := range relation.Arguments {
		newRelation.Arguments = append(newRelation.Arguments, argument.Copy())
	}
	return newRelation
}

// Returns a new relation, that has all variables bound to bindings
func (relation Relation) BindSingle(binding Binding) Relation {

	boundArguments := []Term{}

	for _, argument := range relation.Arguments {
		arg := argument.Bind(binding)
		boundArguments = append(boundArguments, arg)
	}

	return NewRelation(relation.Negate, relation.Predicate, boundArguments)
}

// Returns multiple relations, that has all variables bound to bindings
func (relation Relation) BindMultiple(bindings BindingSet) []Relation {

	boundRelations := []Relation{}

	for _, binding := range bindings.GetAll() {
		boundRelations = append(boundRelations, relation.BindSingle(binding))
	}

	return boundRelations
}

func (relation Relation) IsBound() bool {
	for _, arg := range relation.Arguments {
		if arg.IsVariable() || arg.IsAnonymousVariable() {
			return false
		}
	}

	return true
}

// check if relation uses a variable (perhaps in one of its nested arguments)
func (relation Relation) UsesVariable(variable string) bool {

	var found = false

	for _, argument := range relation.Arguments {
		found = found || argument.UsesVariable(variable)
	}

	return found
}

func (relation Relation) ConvertVariablesToConstants() Relation {

	newArguments := []Term{}

	for _, argument := range relation.Arguments {

		newArgument := argument.ConvertVariablesToConstants()
		newArguments = append(newArguments, newArgument)
	}

	return NewRelation(relation.Negate, relation.Predicate, newArguments)
}

func (relation Relation) ReplaceTerm(from Term, to Term) Relation {
	newRelation := NewRelation(relation.Negate, relation.Predicate, []Term{})
	for _, argument := range relation.Arguments {
		newArgument := argument.ReplaceTerm(from, to)
		newRelation.Arguments = append(newRelation.Arguments, newArgument)
	}
	return newRelation
}

func (relation Relation) String() string {

	args, sep := "", ""

	for _, Argument := range relation.Arguments {

		args += sep + Argument.String()
		sep = ", "
	}

	sign := ""
	if relation.Negate {
		sign = "-"
	}

	return sign + relation.Predicate + "(" + args + ")"
}

func (relation Relation) IndentedString(indent string) string {

	args := ""
	sep := ""

	for _, Argument := range relation.Arguments {

		if Argument.IsRelationSet() {
			args += sep + Argument.TermValueRelationSet.IndentedString(indent+"    ")
		} else {
			args += sep + Argument.String()
		}

		sep = ", "
	}

	sign := ""
	if relation.Negate {
		sign = "-"
	}

	return "\n" + indent + sign + relation.Predicate + "(" + args + ")"
}
