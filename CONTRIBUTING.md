# Contributing

Hello, we are very happy you decided to contribute to go mail. But before you start with your contribution, 
please make sure to read through our guidelines:

- [Issue Reporting Guidline](#issue-reporting-guidline)
- [Pull Request Guidlines](#pull-request-guidlines)

## Issue Reporting Guideline

If you find a bug or believe that some important feature is missing you can open a new issue on the Github-Project page, 
using our provided issue templates. Before creating a new issue, please make sure that there isn't already an issue 
covering this problem or requesting this feature.

## Pull Request Guidelines

- The `main` branch always contains the latest stable released version and doesn't take PRs. 
  Instead, create dedicated feature branches and submit your PR to our `dev` branch.
- It's okay if your PR contains several small commits as we will squash the PR before merging it.
- Please try to use meaningful commit messages.
- Before creating a new PR, check if your code is linted correctly.
- If you want to add a new feature:
    - Add a small but complete description of the new feature.
    - Please provide a convincing reason why you think this feature needs to be added.
- If you add a bug fix:
    - Please refer the corresponding issue, if one exists, in your PR.
    - If no issue exist for the bug you fix you need to provide a detailed description of the error and if possible a live demo. Or create a new issue on our [Github page](https://github.com/ainsleyclark/go-mail/issues)
- Create unit tests for new features and use the `make all` command to test, lint and format.