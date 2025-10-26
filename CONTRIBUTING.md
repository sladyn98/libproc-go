# Contributing to libproc

First off, thank you for considering contributing to libproc! It's people like you that make this library better.

## Code of Conduct

By participating in this project, you are expected to uphold our commitment to providing a welcoming and inspiring community for all.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When you are creating a bug report, please include as many details as possible:

* **Use a clear and descriptive title**
* **Describe the exact steps which reproduce the problem**
* **Provide specific examples to demonstrate the steps**
* **Describe the behavior you observed after following the steps**
* **Explain which behavior you expected to see instead and why**
* **Include Go version, macOS version, and architecture (Intel/Apple Silicon)**

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, please include:

* **Use a clear and descriptive title**
* **Provide a step-by-step description of the suggested enhancement**
* **Provide specific examples to demonstrate the steps**
* **Describe the current behavior and explain which behavior you expected to see instead**
* **Explain why this enhancement would be useful**

### Pull Requests

* Fill in the required template
* Do not include issue numbers in the PR title
* Follow the Go coding style
* Include thoughtfully-worded, well-structured tests
* Document new code
* End all files with a newline

## Development Process

### Setting Up Your Development Environment

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/yourusername/libproc.git
   cd libproc
   ```

3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/original/libproc.git
   ```

4. Create a branch for your changes:
   ```bash
   git checkout -b feature/my-new-feature
   ```

### Making Changes

1. Make your changes in your feature branch
2. Add or update tests as needed
3. Ensure all tests pass:
   ```bash
   go test -v -race ./...
   ```

4. Format your code:
   ```bash
   gofmt -s -w .
   ```

5. Run static analysis:
   ```bash
   go vet ./...
   ```

6. Run linters:
   ```bash
   golangci-lint run
   ```

### Committing Your Changes

* Use clear and meaningful commit messages
* Follow the convention:
  ```
  type: subject

  body (optional)
  ```
  Where type is one of: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

Examples:
```
feat: add support for proc_pidinfo

This adds a new function to retrieve detailed process information
using the proc_pidinfo syscall.
```

```
fix: handle null terminator in PidName correctly

Previously, process names might include trailing null bytes.
This fix properly strips them.
```

### Submitting Your Changes

1. Push your changes to your fork:
   ```bash
   git push origin feature/my-new-feature
   ```

2. Create a Pull Request from your fork to the main repository
3. Wait for the CI checks to pass
4. Address any review comments

## Style Guidelines

### Go Code Style

* Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
* Use `gofmt` to format all code
* Use meaningful variable and function names
* Add comments for exported functions and types
* Keep functions focused and small
* Handle errors explicitly

### Testing Style

* Write table-driven tests where appropriate
* Use meaningful test names: `TestFunctionName_Scenario`
* Include both positive and negative test cases
* Add benchmarks for performance-critical code
* Aim for high code coverage (>80%)

### Documentation Style

* Document all exported functions, types, and constants
* Include usage examples in doc comments
* Use proper Go doc comment format
* Update README.md when adding features

## Project Structure

```
libproc/
├── README.md              # Project overview and usage
├── LICENSE                # BSD-3-Clause license
├── CONTRIBUTING.md        # This file
├── doc.go                 # Package documentation
├── syscall.go             # Main Go syscall wrappers
├── syscall_arm64.s        # ARM64 assembly trampolines
├── tests/                 # Tests directory
│   └── libproc_test.go    # Unit tests
└── .github/
    └── workflows/
        └── ci.yml         # CI/CD configuration
```

## Testing

### Running Tests

```bash
# Run all tests
go test ./tests/

# Run tests with coverage
go test -coverprofile=coverage.out ./tests/
go tool cover -html=coverage.out

# Run tests with race detection
go test -race ./tests/
```

### Writing Tests

All new features should include:

1. Unit tests in `tests/libproc_test.go`

## Questions?

Feel free to open an issue with your question, or reach out to the maintainers directly.

## License

By contributing to libproc, you agree that your contributions will be licensed under the BSD 3-Clause License.
