# Agreement - implementation

Agreement checking is done by the `AgreementChecker`. It uses the `entityTags` from `DialogContext`.

Agreement checking is done in the `check_agreement` step of `respond.rule`, but also during anaphora resolution.