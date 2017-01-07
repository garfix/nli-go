# Syntactic rewrite rules from Natural Language Programming with ThoughtTreasure - Erik T. Mueller (p. 142)

```
E -> A          ADJP -> adjective                       blue
E -> B E        ADJP -> adverb ADJP                     very blue
E -> D E        ADJP -> determiner ADJP                 
E -> E B        ADJP -> ADJP adverb                     
E -> E E        ADJP -> ADJP ADJP                       very blue and bright red
E -> E Y        ADJP -> ADJP PP                         very blue on the edges
E -> K E        ADJP -> conjunction ADJP                and bright red
W -> B W        VP -> adverb VP                         
W -> H W        VP -> pronoun VP                        who like theater
W -> R W        VP -> preposition VP                    
W -> V          VP -> verb                              reads
W -> W 0        VP -> VP expletive                      reads it
W -> W B        VP -> VP adverb                         was easy                                        (overbodig?)
W -> W E        VP -> VP ADJP                           was easy
W -> W H        VP -> VP pronoun                        see you
W -> W V        VP -> VP verb                           
W -> W X        VP -> VP NP                             kicks the ball, understanding her
W -> W Y        VP -> VP PP                             falls on the floor
X -> D X        NP -> determiner NP                     the dress
X -> E X        NP -> ADJP NP                           beautiful red dress
X -> H          NP -> pronoun                           they, whom
X -> K X        NP -> conjunction NP                    and the quick fox
X -> N          NP -> noun                              dog
X -> X 9        NP -> NP element                        Bob 's
X -> X E        NP -> NP ADJP                           la robe bleue                                   (french only)
X -> X W        NP -> NP VP                             the friends who like theater
X -> X X        NP -> NP NP                             the cat and the dog
X -> X Y        NP -> NP PP                             the dog on the floor
X -> X Z        NP -> NP S                              the friends with whom she went to play
X -> Z          NP -> S                          
Y -> B Y        PP -> adverb PP
Y -> K Y        PP -> conjunction PP                    and on the floor
Y -> R B        PP -> preposition adverb
Y -> R X        PP -> preposition NP                    on the mat
Y -> Y Y        PP -> PP PP                             in a bar under the sea
Z -> B Z        S -> adverb S                           carefully he ran upstairs
Z -> E W        S -> ADJP VP
Z -> E Z        S -> ADJP S
Z -> H X        S -> pronoun NP
Z -> K Z        S -> conjunction S                      and the dog ate
Z -> U          S -> interjection
Z -> W          S -> VP                                 wash the dishes!
Z -> x          S -> SententialLexicalEntry
Z -> X E        S -> NP ADJP
Z -> X W        S -> NP VP                              the cat sleeps
Z -> X Z        S -> NP S                               the friends whom she took to the theater, whom she took to the theater
Z -> Y Z        S -> PP S                               with whom she went to play
Z -> Z B        S -> S adverb
Z -> Z Z        S -> S S                                the cat sat and the dog ate
```

Many of these rules are restricted by conditions named in the book.
