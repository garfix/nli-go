package central

type EntityLabel struct {
	label      string
	variable   string
	activation int
}

type EntityLabels struct {
	labels []EntityLabel
}

func NewEntityLabels() *EntityLabels {
	return &EntityLabels{
		labels: []EntityLabel{},
	}
}

func (e *EntityLabels) Clear() {
	e.labels = []EntityLabel{}
}

func (e *EntityLabels) Copy() *EntityLabels {
	newLabels := []EntityLabel{}
	newLabels = append(newLabels, e.labels...)
	return &EntityLabels{
		labels: newLabels,
	}
}

func (e *EntityLabels) GetLabel(label string) (EntityLabel, bool) {
	for _, aLabel := range e.labels {
		if aLabel.label == label {
			return aLabel, true
		}
	}

	return EntityLabel{}, false
}

func (e *EntityLabels) SetLabel(label string, variable string) {
	aLabel, found := e.GetLabel(label)
	if found {
		aLabel.variable = variable
	} else {
		aLabel = EntityLabel{
			label:      label,
			variable:   variable,
			activation: 3,
		}
		e.labels = append(e.labels, aLabel)
	}
}

func (e *EntityLabels) IncreaseActivation(label string) {
	aLabel, found := e.GetLabel(label)
	if found {
		aLabel.activation++
	}
}

func (e *EntityLabels) DecreaseActivation() {
	newLabels := []EntityLabel{}

	for _, aLabel := range e.labels {
		if aLabel.activation > 1 {
			aLabel.activation--
			newLabels = append(newLabels, aLabel)
		}
	}

	e.labels = newLabels
}
