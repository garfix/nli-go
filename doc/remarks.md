# 2018-12-15

Merged the branch askuser into master. All tests passed again. I am happy, this all was a large improvement. Now I have to do some cleanup and
a lot of communication to my users (yes, there aren't any, I know, haha).

# 2018-12-01

I will stop solving relation sets, and in stead solve relation routes and relation groups.

# 2018-11-19

I finally managed to implement asking the user which of the persons named Lord Byron s?he means.

Now I have sense() relations. I think I turn these into a simple array of sense information.

Once I have sense information about the senses of each of the names in a query, I think I must remove the name() relations from the queries I sent to the knowledge bases, and replace the entity-variables with their entity-ids.

For example:

in:

name(E1, N) born(E1, B) has_child(E1, X)

extract senses

born(E1, B) has_child(E1, X)
sense(E1, 'db1', 8311)
sense(E1, 'db2', 136)

query databases

born(E1, B)
    db1: born(E1, B) -->  born(8311, B)  -->  B = '1880-10-09'
    db2: born(E1, B) -->  born(136, B)   -->  B = nil

has_child(E1, X)
    db1: has_child(E1, X)  -->  has_child(8311, X)  -->  X = nil
    db2: has_child(E1, X)  -->  has_child(135, X)   -->  X = 2

result:
    B = '1880-10-09', X = 2

Note: this is the solution of how to find the answer of a question whose information is spread over multiple databases.

# 2018-11-18

Steps for selecting a person from a given name that occurs multiple times.

* User asks a question
* System extracts the names
* System finds database/id identities for each name, along with a sort of description of each identity, for the user
* System stores the original question, along with the open question original_input(), open_question()
* System outputs these as a multiple choice question to the user
* User answers the question with a single answer
* System stores answer_open_question() and removes open_question()
* System finds database/id identities again for each name, along with a sort of description of each identity, for the user
* System selects the ones that match the user's answer from answer_open_question()
* System removes answer_open_question() and adds name_information()

From that point on, name_information() is present for the name and this will be used whenever the name is used in the dialog.

# 2018-11-11

I am dropping the "suggest" option that auto-suggests the next word in a sentence. It didn't really help the user.
Many suggestions were made that were syntactically correct, but not semantically.

The consequence of this is that we must provide example sentences to the user.
Or, more formally, we must specify to the user exactly what the format is that the sentences must adhere to.

This is actually not such a bad idea. Humans always expect an nli to understand everything he asks it.
But that may cease when the user is given a list of allowed sentences and phrases he can use.
I think it is a necessary and reasonable idea in human-computer interaction.

# 2018-10-16

Current form for name recognition: an entities.json file with, per entity, its name, a "knownby" with at least a "name", and possibly some other fields to help identification.

The id of the entity is mapped to the variable Id. The relevant other variable is stored in N.
Each knownby is a relation set (no a relation), consisting of database specific relations.

    {
      "person": {
        "knownby": {
          "name": "foaf_name(Id, Name)",
          "birth_date": "birth_date(Id, Name)",
          "birth_place": "birth_place(Id, Name)"
        }
      }
    }

# 2018-10-06

A name can be an identifier of different things. The name of a person, or the code of a mineral.
Different relations are used to name these things. Different fields are needed to disambiguate between different entities of the same field.

    person: name(), birthdate(), birthplace()
    town: name(), state(), country()

Some but not all of the fields may be present. You need to know the entity's type in order to determine the fields.

# 2018-09-24

I plan to postpone proper name recognition until after the parse phase. But only _directly_ after it.

This means that I will recognize "Lord Byron" as a name of two words. I will not recognize that it it stored in DBPedia at this time.

Until now I recognized names by their capitalized names (L)ord (B)yron and insertions (van, de). I will not do this any more. The main reason is that this principle does not hold for most of the proper names in the rest of the world.

I now plan to recognize _any word_ as a proper noun. The advantage of this is that finding the name in the database can be postponed until after the parse phase and be put in its own phase.

I assumed that I would have to try the longest possible proper noun first, because the match of a longer proper noun is always preferable. However, matching a long proper name built from _any words_ is a recipe for problems. And the longest possible proper name does not exist. It may take up to 10 words or more.

So what I'll do is I'll start with 1 word proper nouns. If it is still possible to parse the whole sentence, the one word noun is fine.

Name recognition just should not be done _before_ the parse phase. Exactly because anything can be a name.

======================

Name resolution phase takes as input generic relations, and adds to these a number of senses:

    name(E5, "John") => name(E5, "John") sense(E5, 'dbpedia', <http://dbpedia.org/resource/John>)

Each store has at most 1 sense, and the "reference" of these senses must be the same (i.e. the same person occurs in multiple databases)

Then, when data store 'dbpedia' is accessed, the variable 'E5' is replaced by the store specific ID <http://dbpedia.org/resource/John>

# 2018-09-23

I will do it a little bit differently. The Dialog Context is only contacted by the system, and the system places the answer in the DC.

# 2018-09-20

Yes I am still thinking about how to ask the user something.

But I have decided not to break the application up into tasks. It's a pretty idea, but would require forcing it into strange unnatural structures, which isn't good.

I will continue with an alternative that I will now explain:

The system uses a Dialog Context, which is persisted somewhere where other processes can access it too.

When the system has a question it
* checks the Dialog Context for an answer, or
* stops the application with a question for the user

If the application is stopped, it is up to the process that ran the application to ask the user the question, and place the answer in the Dialog Context.

Next, the same question is passed to the system. The system processes the question in the same way. Except when it checks the Dialog Context, it succeeds and continues.

# 2018-09-13

It occurred to me that the technique of using tasks is also a solution for a problem that I had stuck away.

Syntactic and semantic analysis of a sentence can be ambiguous: it can result in multiple possibilities. Up until now I have chosen to pick only the first one of these. I simply had no way of dealing with any others and I expected there were not many cases in with this mattered.

However, when the process is split up in tasks, the task of parsing may yield 2 or three parse trees, and these can all be dealt with.

In the days before I have found:

* the task T
* the task sequence [T T T]

I now add to this

* the task switch (T | T | T)

The result of the task switch is the first result, or null if none yields a result. When one path has a result, the others need not be followed.

# 2018-09-02

Let's see what happens when there are 2 databases and both have Lord Byron.

Database 1:

* 16882 name "Lord Byron"

Database 2:

* <http://name.org/byron> name "Lord Byron"
* <http://name.org/byron_(umpire)> name "Lord Byron"

Preprocessing:

Lord Byron:
    * db 1: id = 16882, description = "English Lord that lived from ..."
    * db 2: id = <http://name.org/byron>, description = "English Lord that lived from ..."
    * db 2: id = <http://name.org/byron_(umpire)>, description = "American baseball umpire"

Input: "Name Lord Byron's children"

Name resolution phase: in this imperative sentence the NP can be resolved into the proper name "Lord Byron". The system then knows of three persons Lord Byron in two databases. It must determine which Lord Byron the user means.

This could be a way:

"By 'Lord Byron' did you mean":

Database 1:
(o) English Lord that lived from ...
( ) None of the above

Database 2:
(o) English Lord that lived from ...
( ) American baseball umpire
( ) None of the above

"Please choose one option per database, and press OK"

And the user is given two radio groups, one for each database. Multiple answers may be checked, but only one per database.

Once the choice is made, the meaning of the sentence could contain some database specific information:

"Lord Byron" => name(A, "Lord Byron"), sense(A, 16882, person, db1), sense(A, person, <http://name.org/byron>, db2)

The meaning of "sense" is simply that of the logical sense: one of many representation of a single meaning

https://en.wikipedia.org/wiki/Sense_and_reference

This is actually much better then what I have been doing so far, where the variable A is given the id of the database. The new form is better suited to integrate information of two or more databases.

Since the system is not able to automatically determine if the same name in two database refers to the same person in real life, it simply asks the user. This needs to be done only once per session.

# 2018-08-31

I picked up where I left off. There are two persons called "Lord Byron" in DBPedia. One of them is a baseball umpire.
So when someone asks a question about Lord Byron, the system needs to ask "Which one?" followed by some meaningful description of the two men that sets them apart in a way that is easily distinguishable by the user.

If you think this is a highly improbable case, check "Michael Jackson". However, in both cases, the first person with the same name is more important or primary than the others. Nevertheless, it should at least be possible to query the other ones and they should not be mixed up.

This case is problematic for the following reasons:

* The system now needs to enter a dialog with the user (no dialog was used so far)
* This dialog needs to be independent of the user interface. Cannot just do a "readline", the user interface may be completely different. It may be a web interface for instance.
* Also, the process of answering the question may need to be stopped and picked up later, when the user has made his choice.
* The options that follow "Which one?" are database specific.
* The name "Lord Byron" may occur in multiple databases that are used in the session with the user. For some questions it is necessary to know which Lord Byron from the first database corresponds to the Lord Byron of the second database.
* In what part of the process are the names recognized and in what part does the user choose one over the other?

A _name_ is a special form in syntactic processing, because proper names are not in the dictionary. Currently I am handling this by requiring that names are written with capitals and only certain combinations of capitalized words and insertions are allowed. This is far from ideal, especially as we start to extend proper names to other things than persons.

For this reason it would be good to search the database for names that might match "proper noun" parts. An idea is to add all proper names to the dictionary, but this requires regular syncs from the database to the dictionary, and these syncs may be quite large.

However, if one wants to support correction of mis-spellings, preprocessed names may be necessary. ("Lord Biron? Did you mean Lord Byron?")

Using database specific structures in the parsing phase contrasts with the design I used to far, that database specific data is used only in specific places of the answering phase. However, there is no good reason to keep it restricted to that place. Using databases forces one to make choices that are enforced by the database.

# 2018-01-27

Working on the question "Who was X's father?". The sentence has the structure 'isa(Z, who) identity(Z, Y) father(X, Y)', but the condition is now

    question(S, wh_question) isa(X, who) has_father(A, B)

I left out the 'identity' clause. The reason for this is that both 'isa(Z, who) identity(Z, Y)' cannot match anything on their own.
They could be used as the final relations in a set, in which case the 'identity' relation would unify the variables Z and Y. But for the present example
there is no use for it, and I will handle such a case only when it presents itself.

# 2018-01-21

I have done what I described. A problem can now have multiple solutions.
I added the possibility to transform an input sentence so that it forms a different type of solution.
I also changed the place where the scoping is performed. It used to be a major step in the pipeline,
    now it is only executed just before solving the problem, and after the solution transformations have been executed.

# 2018-01-19

DBPedia has two ways of storing parent-child relationships. It has the dbo:child relation and the dbp:children relation. The former stores links to other persons (children),
the latter stores the number of children. When the first is given, the second is not (I think) and vice versa.

So when I want to know how many children a person has, I need to try both ways. The standard query that is produced from the input sentence maps to the dbo:child relation.
If the standard query fails, I want to be able to try the other way. This means that I want to try 2 solutions for the same problem.

Currently I am trying only a single solution and when that fails, the answer is that it is unknown how many children the person has.

I can easily change this into trying all solutions. There are some issues now:

- each solution has a "no results" section. Which section should I use if all solutions fail?
- or should I do a single solution with multiple sub-solutions?
- are different solution paths database-dependent (db layer) or solution dependent (domain specific). I tend to go for the last option.
- can I make the solution "rephrase the question" in order to answer it?

If I would rephrase the question about the number of children, this would become something like this:

    have_child(A, B) how(_) many(_) => have_n_children(A, N)

A plain verb is turned into an abstract, problem specific, relation. Is this still the domain model, or is it the database model?
Is it even possible to solve this at the database level? Must it be solved at the database level?

# 2018-01-14

I released http://patrickvanbergen.com/dbpedia/app/ and mentioned it on Twitter. It is very unimpressive for a demo, but I had to release _something_ so that I can show 
other people what I am doing. It will be nice to gradually improve upon it. 

# 2018-01-09

This picture is just to clarify things to myself. A, B, C, and D are range variables in nested scopes.

A1 B1 C1 D1
A1 B1 C1 D2
A1 B1 C2 D1
A1 B1 C2 D3
A1 B2 C1 D1
A1 B2 C1 D4
A1 B2 C2 D1
A1 B2 C2 D5

When evaluating scopes, each scope is evaluated only once, to restrict the number of queries.
As you can see D1 and D2 occur multiple times in the result set. 
Suppose the quantifiers are: B = 2, C = 2, D = 2
This can be implemented as follows:
At quant C, all quants in C's scope are evaluated, which means D.
For all unique values C, the number of distinct values D should equal the quantifier of D, i.e. 2.
=> This does not hold, since below C2 are D2, D3, D4, and D5.
This means that the scopes can only be evaluated after the outermost scopes are evaluated, which means: after the query is done.

After the outermost relation set has been solved, the quantifiers of its quants should be evaluated. This affects the resulting bindings.

I think I am going to go with a special construct in SolveQuant. It keeps track of a global quant level variable.
This is not pretty, but it works. When the quant is done, and this is the outermost level (0), the quantifiers are evaluated.
I will chose this solution, because the function solveQuant() will solve the quant _fully_ and will not be dependent on some
outside function call. Furthermore, this attempt may not even succeed, and I don't want to rewrite my code too much for a failed
attempt. I will place a todo to be resolved later, if the attempt should succeed.

# 2018-01-07

I have done it that way, but I have introduced an IF / THEN construct 

    IF isa(E2, child) THEN subject(S5, E1) object(S5, E2) isa(S5, have) => has_child(E1, E2)
    
because this is much clearer (I think) and I didn't want to reuse any of the brackets for this purpose, because they already signify other things.

I am now running into the problem that in scoping, the range, which is always evaluated first, yields too many values.

If the range is isa(E, child), E resolves to all persons (who are children). This is too much for any real-world database.

So I am now thinking about combining the Range with the Scope. But i wonder why I didn't think of that before. I knew it crossed my mind. What was the reason?

---

I think this is it. A quant has three aspects, a Range, a Quantifier and a Scope. For example: dogs, 3 or more, have a bone. The relation between these aspects is thus:

The range (dogs) is the domain of discourse. And if that was all it could be just a part of the Scope ( => dog(X) has_bone(X) ), which limits the domain as well.

The quantifier is a check. Does the scope yield 3 entities? Does the scope variable yield 3 possible values?

The range then, is not necessary in most cases. To check if the scope yields 3 entities, it is not necessary to know that these entities are dogs. If they were cats, 
it would be just as well.

But when the Quantifier is _all_, or _most_, the scope becomes important. Because the number of dogs is different from that of cats. 
When the Quantifier is _all_, the number of results in the scope must match the number of dogs. COUNT(DISTINCT S1) = COUNT(R) where R is just the query isa(R, dog) 

The number of different values for R may be large, and even very large, but its number just needs to be counted. The individual values are not used for further processing.

The range relation set may be copied to extend the scope relation set, but it must remain independent as well.  

# 2018-01-03

I think I have found a better idea. Neither of the proposed alternatives from yesterday where really great. Also, allowing relation rewrites
without minding scope borders seems like a bad way to go.

What I am now thinking about is this

    [isa(E2, child)] subject(S5, E1) object(S5, E2) isa(S5, have) => has_child(E1, E2)

Which means: rewrite 'subject(S5, E1) object(S5, E2) isa(S5, have)' to 'has_child(E1, E2)',
if 'isa(E2, child)' occurs in the sentence. 'isa(E2, child)' is unaltered by the rewrite and may live in any scope of the sentence.

# 2018-01-02

Hi there! Happy 2018! I used to scope Quantifier and Range relations in the Relationize phase, quite early in the process.
But the problem was that it got in the way of generic 2 domain specific conversions, they became too complex.
Doing the scoping later (in the scoping phase where the quant is formed) is problematic too because collecting all relations
that contain the range variable results in too many relations. This is because the variable is used in other relations
higher up the parse tree as well. I will try to visualize:

    np(E1)       posession(E1, E2)
    |
    nbar(E1)
    |
    dp(D1) nbar(E1)

D1 forms the quantifier; E1 forms the range. The quantifier variable is fine. The range variable is not only bound to relations
below the nbar, but above it as well (in the example: possession(E1, E2))

It is hard to figure out which relations belong to the range and which relations do not. The solution I am now using is based on
the heuristic that each nbar relation forms a specialization. So I will only use the variable on the right side of the specification()

    specification(E1, X)

and all relations that are connected, directly or indirectly to this variable.
