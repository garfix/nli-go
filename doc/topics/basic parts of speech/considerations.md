# Considerations

## Why not use a lexicon?

NLI-GO only has a grammar, it doesn't have a lexicon. Wouldn't it be useful to have a lexicon, with only a single entry for each verb, and not a separate grammar rule for each inflection? The choise for lexicon-as-grammar was made out of flexibility. Language contains many free-form expressions that are not easily caught by a lexicon.

    How many countries have population above 10 million?

"to have population" is not a good entry for a lexicon. It neither belongs to "have", nor to "population". There might be an entry for "have-population" but it would require very specific annotations in how to use it.

The disadvantage of not having a lexicon is that each verb inflection, and each noun inflection needs to be entered separately.
