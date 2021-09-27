# gitlab-multirepo-deployer
Application to trigger multiple gitlab pipelines on a specific branch.  Also able to scan through listed projects to determine if a specified branch is present.

Note - gitlab only allows pipelines to be triggered using CI tokens or trigger tokens - NOT personal tokens.  For manual use, a personal access token is used for general api queries and a trigger token is required per project to trigger pipelines.  When used in CI, the CI job token can be used for queries and triggers.
