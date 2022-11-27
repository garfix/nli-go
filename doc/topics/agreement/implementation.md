# Agreement - implementation

Agreement checking is done by the `AgreementChecker`. It uses the `entityTags` from `DialogContext`.

Agreement checking is done in the `checkAgreement` step. When the step fails, the processing alternative fails.

It is also used during anaphora resolution.

