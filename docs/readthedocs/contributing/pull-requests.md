# Pull Requests

Thank you for taking interest in contributing to CV Wonder!

1. Every Pull Request should have an associated GitHub issue link in the PR description. Please refer to the issue and discussion pages for explanation about this process. If you think your change is trivial enough, you can skip the issue and instead add justification and explanation in the PR description.
1. Your PR is more likely to be accepted if it focuses on just one change.
1. There's no need to add or tag reviewers.
1. If a reviewer commented on your code or asked for changes, please remember to respond with comment. Do not mark discussion as resolved. It's up to reviewer to mark it resolved (in case if suggested fix addresses problem properly). PRs with unresolved issues should not be merged (even if the comment is unclear or requires no action from your side).
1. Please include a comment with the results before and after your change.
1. Your PR is more likely to be accepted if it includes tests.

## Pull Request Title

The pull request title must follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/){:target="_blank"} format.

Format of the title:

```raw
<type>(<scope>): <subject>
```

The type and scope should always be lowercase as shown below.

**Allowed `<type>` values**:

* `feat` for a new feature for the user, not a new feature for build script. Such commit will trigger a release bumping a MINOR version.
* `fix` for a bug fix for the user, not a fix to a build script. Such commit will trigger a release bumping a PATCH version.
* `perf` for performance improvements. Such commit will trigger a release bumping a PATCH version.
* `docs` for changes to the documentation.
* `style` for formatting changes, missing semicolons, etc.
* `refactor` for refactoring production code, e.g. renaming a variable.
* `test` for adding missing tests, refactoring tests; no production code change.
* `build` for updating build configuration, development tools or other changes irrelevant to the user.
* `chore` for updates that do not apply to the above, such as dependency updates.
* `ci` for changes to CI configuration files and scripts
* `revert` for revert to a previous commit

**Allowed `<scope>` values**:

Scope is not restricted to a predefined list of values. It must reflect the part of the codebase that is being changed.

**Breaking changes**:

A Pull Request introducing a breaking API change needs to append a ! after the type/scope.

### Example titles

```raw
feat(ui): add graph in claim view
feat(api)!: delete the existing CLI flag
fix(api): handle dependencies in the claim view
docs: add the contribution guide
```

The type/scope `chore(deps)` is generally used by Dependabot for dependencies updates. You can also use this pattern if needed. Dependencies updates should be contained in a single Pull Request to ensure that the changes are properly tested and reviewed.

```raw
chore(deps): bump gorm.io/gorm from 1.25.4 to 1.25.12
```

## Commits

Each commit message doesn't have to follow the conventions as long as it is clear and descriptive since it will be squashed and merged.
