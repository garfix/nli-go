2017-01-07

I decided to work with releases. Each release has a goal functionality, and must be documented so as to be usable to others.

I cannot just use Erik T. Mueller's syntax rules (mueller-rewrites), because they have many constraints. I prefer to solve these constraints in the rules themselves (if that's possible). I keep them for inspiration.

2016

Als je wilt dat de representatie een Horn clause repr is, moet je NOT en OR expliciet noemen.
Maar is het wel mogelijk om deze in de eerste parse op te nemen, of zijn
Maar is het wel nodig om ze op te nemen? Je kunt de meeste determiners ook niet opnemen.
En je neemt modale elementen (ik dacht dat ..) ook niet op
Ok, maar daarmee is je representatie echt NIET logisch te noemen
Niet alleen komen EN en OF niet overeen met hun logische equivalenten en is pragmatische interpretatie mogelijk,
    ook is keihard weergegeven NIET te beperkt, omdat er ook MISSCHIEN en NAUWELIJKS bestaan.
FOPC is in zijn algemeenheid gewoon te beperkt, en er is geen goed alternatief.

Er is een probleem met left-recursion in de simpele parser NP :- NP VP

een agent

agent: {
    grammar: {
        rules ...
    }
    lexicon: [
        entries ...
    ]
}

een lexicon op zichzelf:

lexicon: [
    {
        form: ..
        pos: ..
        sense: ..
    }
    {
        form: ..
        pos: ..
        sense: ..
    }
]




		predication(S1, marry)
		object(S1, E2)
		subject(S1, who)
		name(E1, 'Kurt Cobain')

Ik maak 'grammatical_subject' nu het predicaat dat aangeeft waar de hoofdzin is. Een predicatie-object is niet aanwezig in het domein-specifieke resultaat. 
Dit grammatical_subject geeft ook aan of de zin actief is of passief.
