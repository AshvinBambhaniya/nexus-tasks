# Contributing to Nexus Tasks

First off, thanks for taking the time to contribute! 🎉

The following is a set of guidelines for contributing to Nexus Tasks. These are mostly guidelines, not rules. Use your best judgment and feel free to propose changes to this document in a pull request.

## Code of Conduct

This project and everyone participating in it is governed by the [Nexus Tasks Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## How Can I Contribute?

### Reporting Bugs

This section guides you through submitting a bug report for Nexus Tasks. Following these guidelines helps maintainers and the community understand your report, reproduce the behavior, and find related reports.

- **Use a clear and descriptive title** for the issue to identify the problem.
- **Describe the exact steps which reproduce the problem** in as much detail as possible.
- **Provide specific examples** to demonstrate the steps.
- **Describe the behavior you observed** after following the steps and point out what exactly is the problem with that behavior.
- **Explain which behavior you expected to see instead and why.**

### Suggesting Enhancements

This section guides you through submitting an enhancement suggestion for Nexus Tasks, including completely new features and minor improvements to existing functionality.

- **Use a clear and descriptive title** for the issue to identify the suggestion.
- **Provide a step-by-step description of the suggested enhancement** in as much detail as possible.
- **Explain why this enhancement would be useful** to most Nexus Tasks users.

### Pull Requests

1. Fork the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. If you've changed APIs, update the documentation.
4. Ensure the test suite passes.
5. Make sure your code lints.
6. Issue that pull request!

### Pre-commit Hooks

We use [pre-commit](https://pre-commit.com/) to ensure code quality (formatting, linting, security) and commit message compliance.

1.  **Install pre-commit:**

    ```bash
    pip install pre-commit
    # Or via brew: brew install pre-commit
    ```

2.  **Install the hooks:**

    ```bash
    pre-commit install
    pre-commit install --hook-type commit-msg
    ```

3.  **Run manually (optional):**
    ```bash
    pre-commit run --all-files
    ```

## Development Setup

### Backend (Go)

1.  Navigate to the backend directory:
    ```bash
    cd backend
    ```
2.  Install dependencies:
    ```bash
    go mod download
    ```
3.  Run the API server:
    ```bash
    go run cli/main.go api
    ```

### Frontend (Nuxt 4)

1.  Navigate to the frontend directory:
    ```bash
    cd frontend
    ```
2.  Install dependencies:
    ```bash
    npm install
    ```
3.  Run the development server:
    ```bash
    npm run dev
    ```

## Styleguides

### Git Commit Messages

We follow [Conventional Commits](https://www.conventionalcommits.org/). This is enforced by `commitlint`.
Example: `feat: add task assignment` or `fix: resolve websocket reconnection issue`.

### Coding Standards

- **Go:** Follow standard `gofmt` and `golangci-lint` recommendations.
- **Vue/TypeScript:** Follow Nuxt 4 best practices and ESLint configurations.
