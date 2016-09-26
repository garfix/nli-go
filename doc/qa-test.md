
goal: married_to(A, B), gender(A, G), name(B, N)

goal: married_to(A, B), gender(A, G), name(B, N), name(A, 'Kurt Cobain')

		married_to(A, B) :- marriages(A, B)
		name(A, N) :- person(A, N, _, _)
		gender(A, male) :- person(A, _, 'm', _)
		gender(A, female) :- person(A, _, 'f', _)

		married_to(A, B), name(A, AN), name(B, BN) :- marriages(AN, BN)
		name(A, N) :- person(A, N, _, _)
		gender(A, male) :- person(A, _, 'm', _)
		gender(A, female) :- person(A, _, 'f', _)

hypothese 1
een entity-variabele wordt niet gematched aan een integer of string, alleen aan een variabele (abstract)
alleen zo kun je meerdere databases gebruiken

hypothese 2
pas eerst een transformatie toe van het hele doel naar de db-taal
probeer daarna pas het doel op te lossen
doe dit voor elk van de db's

goal in db-taal
marriages('Kurt Cobain', B)
person('Kurt Cobain', _, _)
person(B, G, _)

zullen we de db dit gewoon laten oplossen? het probleem hiermee is dat je er geen domein-specifieke regels op los kunt laten

