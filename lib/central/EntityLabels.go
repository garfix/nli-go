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
