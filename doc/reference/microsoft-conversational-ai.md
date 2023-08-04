Some interesting features from Microsoft's Conversational AI Platform

Book: Microsoft Conversational AI Platform for Developers - Stephan Bisser

## Intent recognition

Intents are recognized by trigger sentences like

    Name of trigger: BuySurface
    Trigger phrases:
        How can I buy {ProductType = Surface PRO}
        I want to buy {ProductType = Surface PRO}

Note: multiple trigger phrases are attached to the same trigger

## Memory scopes

The platform uses multiple scopes:

* user scope: data scoped to the ID of the user you are conversing with
* conversation scope: data scoped to the ID of the conversation you are having
* dialog scope: data for the life of the associated dialog, providing memory space for each dialog to have internal persistent bookkeeping. Dialog scope is cleared when the associated dialog ends
* turn scope: data that is only scoped for the current turn. The turn scope provides a place to share data for the lifetime of the current turn
* settings scope: any settings that are made available to the bot via the platform-specific settings configuration system
* this scope: the active action’s property bag. This is helpful for input actions since their life type typically lasts beyond a single turn of the conversation
* class scope: instance properties of the active dialog


