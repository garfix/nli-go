# Syntactic rewrite rules from The structure of modern english - Laurel J. Brinton

```
Det -> Art                                          a, an, the
Det -> Dem                                          this, that, these, those
Det -> Poss                                         that man 's, my, our, their
Det -> Q                                            some, any, every, each, neither, more
Det -> Wh-word                                      which, what, whose

    Deg                                             more, most, less, least, very, quite                (degree adverb)
    PSpec                                           right, straight, one mile, three seconds            (prepositional specifiers)

adjective -> adjective conjunction adjective        long and boring

ADJP -> adjective                                   fierce
ADJP -> Deg adjective                               very fierce
ADJP -> ADVP adjective
ADJP -> adjective PP
ADJP -> Deg adjective PP
ADJP -> ADVP adjective PP
ADJP -> ADJP conjunction ADJP                       very slow and quite tedious

adverb -> adverb conjunction adverb                 quietly and smoothly

ADVP -> adverb
ADVP -> Deg adverb
ADVP -> ADVP conjunction ADVP                       very cautiously but quite happily

noun => noun conjunction noun                       cats and dogs

NP -> NBar                                          large dog on the sofa
NP -> Det NBar                                      the large dog on the sofa
NP -> pronoun                                       he
NP -> propernoun                                    Goldy
NP -> NP conjunction NP                             the tortoise and the hare

NBar -> noun                                        dog
NBar -> ADJP NBar                                   large dogs, loudly barking dogs
NBar -> NBar PP                                     dog on the sofa
NBar -> NBar conjunction NBar                       cold coffee and warm beer

Poss -> NP possesive-marker                         that man 's                                     ('s must be separated from man)

preposition => preposition conjunction preposition  over or under the covers

PP -> preposition NP                                behind the door
PP -> preposition PP                                from behind the door
PP -> PSpec preposition NP                          right on time
PP -> PSpec preposition PP
PP -> PP conjunction PP                             on the table and under the chair

Aux -> T (M) (Perf) (Prog) (Pass)                                                                   (much work to do here)

Vgp -> verb
Vgp -> verb particle
Vgp -> Aux verb
Vgp -> Aux verb particle

VBar -> Vgp                                         look
VBar -> Vgp NP                                      open a package
VBar -> Vgp NP NP                                   write a friend a letter
VBar -> Vgp NP PP                                   give an excuse to the teacher
VBar -> Vgp ADVP                                    feel lonely
VBar -> Vgp NP ADVP                                 make the dog angry
VBar -> Vgp PP                                      jump into the pool
VBar -> Vgp PP PP                                   talk about the problem with a friend
VBar -> Vgp NP particle                             look it up

VP -> VBar
VP -> VBar PP
VP -> VBar ADVP
VP -> VBar NP

S -> NP VP
S -> NP Aux VP                                      The dog will bite
S -> Aux NP VP                                      Was Zelda leaving for Paris
S -> ADVP S                                         seriously, that's the most shocking news I've heard
S -> PP S                                           to my regret, i have never met Gerard
S -> S ADVP                                         that's the most shocking news I've heard, frankly
S -> S PP                                           i have never met Gerard, to my regret

```
