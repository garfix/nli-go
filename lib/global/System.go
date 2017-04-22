package global

import (
    "nli-go/lib/parse"
    "nli-go/lib/central"
    "nli-go/lib/mentalese"
    "nli-go/lib/parse/earley"
    "nli-go/lib/generate"
)

type system struct {
    log *systemLog
    lexicon *parse.Lexicon
    grammar *parse.Grammar
    generationLexicon *generate.GenerationLexicon
    generationGrammar *generate.GenerationGrammar
    tokenizer *parse.Tokenizer
    parser *earley.Parser
    quantifierScoper mentalese.QuantifierScoper
    relationizer earley.Relationizer
    transformer *mentalese.RelationTransformer
    answerer *central.Answerer
    generator *generate.Generator
    surfacer *generate.SurfaceRepresentation
    generic2ds []mentalese.RelationTransformation
    ds2generic []mentalese.RelationTransformation
}

func (system *system) ImportLexicon(fromLexicon *parse.Lexicon) {
    system.lexicon.ImportFrom(fromLexicon)
}

func (system *system) ImportGrammar(fromGrammar *parse.Grammar) {
    system.grammar.ImportFrom(fromGrammar)
}

func (system *system) ImportGenerationLexicon(fromLexicon *generate.GenerationLexicon) {
    system.generationLexicon.ImportFrom(fromLexicon)
}

func (system *system) ImportGenerationGrammar(fromGrammar *generate.GenerationGrammar) {
    system.generationGrammar.ImportFrom(fromGrammar)
}

func (system *system) Process(input string) (string, bool) {

    tokens := system.tokenizer.Process(input)
    parseTree, _ := system.parser.Parse(tokens)
    rawRelations := system.relationizer.Relationize(parseTree)
    genericRelations := system.transformer.Replace(system.generic2ds, rawRelations)
    domainSpecificSense := system.quantifierScoper.Scope(genericRelations)
    dsAnswer := system.answerer.Answer(domainSpecificSense)
    genericAnswer := system.transformer.Replace(system.ds2generic, dsAnswer)
    answerWords := system.generator.Generate(genericAnswer)
    answer := system.surfacer.Create(answerWords)

    return answer, true
}
