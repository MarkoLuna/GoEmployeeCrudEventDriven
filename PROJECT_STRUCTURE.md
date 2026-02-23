# Standard Go Project Layout

The structure of a Go project varies based on complexity, but common conventions help organize code effectively, primarily using the
cmd/, internal/, and pkg/ directories. 

## Core Directories
A typical, production-ready Go project structure often looks like this:

```bash
project-root/
├── go.mod
├── cmd/
│   ├── api-server/
│   │   └── main.go
│   └── cli-tool/
│       └── main.go
├── internal/
│   ├── auth/
│   ├── handlers/
│   ├── services/
│   └── models/
├── pkg/
│   ├── logger/
│   └── utils/
├── configs/
├── scripts/
└── Makefile
```
 

- `cmd/`: Contains the main entry points for your application(s). Each sub-directory here holds a separate executable program, with its own main.go file declaring package main.
- `internal/`: This directory is for private application and package code that should not be imported by other projects or modules. Go enforces this: the compiler disallows imports of any path containing the element "internal" if the importing code is outside the current module's tree. This is where most of your core business logic, handlers, and service implementations reside.
- `pkg/`: This directory contains public, reusable packages that can be safely used by external applications/projects. If you intend for code to be imported by other Go modules, put it here; otherwise, use `internal/`.
- `go.mod`: The module definition file, located in the project root, which defines the module name and manages dependencies. 

## Other Common Directories

- `configs/`: Stores configuration files or configuration loading logic.
- `scripts/`: Houses various scripts for building, deploying, or running project tasks (e.g., shell, Python scripts).
- `tests/`: Contains extra test data or large integration tests, though most tests (*_test.go) usually live alongside the code they test.
- `docs/`: Project documentation.
- `Makefile`: Automates common developer tasks like building, running tests, or deployments with simple commands (e.g., make build, make test). 

## Key Principles

- Simplicity first: For small projects or learning, a flat structure with a single `main.go` file is often sufficient. Don't over-engineer the structure initially; let it evolve with the project's complexity.
- Avoid `src/`: Unlike some other languages (like Java), Go projects generally do not use a top-level `src` directory for all source code.
- Organize by concern: Within `internal/` or `pkg/`, organize code into logical packages (directories) based on their functionality (e.g., `auth`, `database`, `metrics`).
- Interfaces for decoupling: Use Go's interfaces and dependency injection to keep packages decoupled, which improves testability and maintainability, especially in larger applications. 

For detailed, official guidance, refer to the [Organizing a Go module](https://go.dev/doc/modules/layout) documentation on the official Go website. 

## Other sources
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Organize Like a Pro: A Simple Guide to Go Project Folder Structures](https://medium.com/@smart_byte_labs/organize-like-a-pro-a-simple-guide-to-go-project-folder-structures-e85e9c1769c2)