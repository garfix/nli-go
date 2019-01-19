package knowledge

type sparqlResponse struct {
	Results sparqlResponseResults
}

type sparqlResponseResults struct {
	Bindings []sparqlResponseBinding
}

type sparqlResponseBinding struct {
	Variable1 sparqlResponseVariable
	Variable2 sparqlResponseVariable
}

type sparqlResponseVariable struct {
	Type string
	Value string
	Lang string `json:"xml:lang"`
}
