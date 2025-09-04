# Contributing to ContainDB

First off, thank you for considering contributing to ContainDB! It's people like you that make ContainDB such a great tool. This document provides guidelines and steps for contributing.

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct. Please read [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) to understand what behaviors will and will not be tolerated.

## How Can I Contribute?

### Reporting Bugs

This section guides you through submitting a bug report for ContainDB.

Before submitting a bug report:

- Check the [documentation](README.md) to see if there's a solution to your problem.
- Check if the issue has already been reported in our [Issues](https://github.com/Nexoral/ContainDB/issues) section.

When submitting a bug report:

- Use our bug report template.
- Use a clear and descriptive title.
- Describe the exact steps to reproduce the problem.
- Explain the behavior you expected and what you actually observed.
- Include details about your environment (OS, ContainDB version, Docker version).
- Include screenshots or terminal output if possible.

### Suggesting Enhancements

This section guides you through submitting an enhancement suggestion for ContainDB.

When submitting an enhancement suggestion:

- Use our feature request template.
- Use a clear and descriptive title.
- Provide a step-by-step description of the suggested enhancement.
- Explain why this enhancement would be useful to ContainDB users.
- Include any relevant examples or mockups if applicable.

### Your First Code Contribution

Unsure where to begin? Look for issues labeled with:

- `good-first-issue`: Issues suitable for newcomers.
- `help-wanted`: Issues that need assistance.
- `documentation`: Improvements or additions to documentation.

### Pull Requests

Here's the process for submitting a pull request:

1. Fork the repository.
2. Create a new branch from `main` for your changes.
3. Make your changes in your branch.
4. Ensure your code follows the existing style.
5. Run tests if available.
6. Submit your pull request with a clear description of the changes.

## Development Setup

### Prerequisites

- Go 1.18 or higher
- Docker and Docker Compose
- Git

### Setup Steps

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR-USERNAME/ContainDB.git
   cd ContainDB
   ```
3. Build the binary:
   ```bash
   ./Scripts/BinBuilder.sh
   ```

## Styleguides

### Git Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line
- Consider starting the commit message with an applicable emoji:
  - ✨ `:sparkles:` when adding a new feature
  - 🐛 `:bug:` when fixing a bug
  - 📝 `:memo:` when adding or updating documentation
  - 🚀 `:rocket:` when improving performance
  - 🔧 `:wrench:` when updating configurations
  - ♻️ `:recycle:` when refactoring code

### Go Styleguide

Follow the standard [Go style guide](https://golang.org/doc/effective_go) and use `gofmt` to format your code.

## Additional Notes

### Issue and Pull Request Labels

This section lists the labels we use to track and manage issues and pull requests.

- `bug`: Confirmed bugs or reports that are very likely to be bugs.
- `enhancement`: Feature requests or improvements to existing functionality.
- `documentation`: Improvements or additions to documentation.
- `good-first-issue`: Good for newcomers.
- `help-wanted`: Extra attention is needed.
- `question`: Further information is requested.

## Thank You!

Again, thank you for your contributions. Your efforts help make ContainDB better for everyone!
