# Goals and Principles

The goal of this engine is to provide a means to interact with any database through natural language. Key goals are:

* the system should able to handle complex natural language constructs. 
* the system should be easy to understand. While "easy" may have to be stretched for this domain, any part of it should be extendable with only a few minutes of explanation.
* it should support multiple domains in a single session. The system should able to reason about people and factory parts, if both domains are modelled.
* it should be able to handle multiple languages. It is okay if a change of session is required to switch languages.

## Users

Since this system does not work out-of-the-box for new domains, it needs the following users.

* end user: the one who will use the system. He needs to have knowledge of the domain he wants to ask questions about (of course).
* domain expert: needs to have knowledge of the domain (of course)
* knowledge engineer: queries the domain expert and creates the declarative files
* application programmer: creates a custom program using system classes and declarative files. Knowledge of Go programming language. Some knowledge of the system.
* system programmer: extends the system (so that's me, for now). Needs knowledge of NLP (much), Go (some), Prolog (some), OOP.

## Principles

* ease of use is paramount, for all types of users, but the knowledge engineer has priority
* small queries: the interface to the knowledge base should be as small as possible
