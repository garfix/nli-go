# Relations

I distinguish several stages in processing a sentence.

## Tokenizer

The tokenizer turns a string into an array of strings: tokens. It's kept simple. It doesn't remove "meaningless characters" except for whitespace.
Letters and digits are kept together in a single token. Other character each get their own token.

## Parser

The parser turns a sequence of tokens into a "sense", a set of relations. The Earley parser produces a syntactic parse tree as well, but this is not used any further.

The sense is a literal interpretation of the sentence. It does not really represent the meaning of the sentence. Just see it as the rawest possible relational representation.

It has the form of a set of relations. These standard grammar is quite verbose. It is meant to be complete. You should be able to express any sentence with it.

There's only a limited number of relations. Relations are based on processing function, not on content. Not 'hold(X, Y)' but 'isa(E, hold), subject(E, X) object(E, Y)'

## Generic to Domain Specific transformation

From here on it's relations all the way down. In the next stage the raw, generic, relations are transformed into a domain specific set.

The domain specific layer (DS) lies between the generic layer and the database layer. It has its own representation.

DS is different from generic because generic is too verbose, and the literal interpretation is not the real meaning of the sentence. DS is the real meaning. It expresses the sentence in the relations of the domain.

DS is also different from some database layer, that is likely also relational, because we need to be able to reason about the sentence. This is best done on a level apart from the database, because this is where the concepts are expressed most naturally. Reasoning about siblings is easier with a relation sibling that with a set of parent relations. Also, much information is only implicit in the database. In this layer we can make it explicit. Finally, there's the point of single responsibility. The database relations change when the database structure changes. The DS layer would not change with it.

Each domain has its own relations. The restrictions of the generic form do not hold. holds(X, Y) is fine here.
