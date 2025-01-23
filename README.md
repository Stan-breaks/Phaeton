# Phaeton Language

_A simple interpreted language for modeling stage-based logic and conditional workflows_

---

## Table of Contents

1. [Overview](#overview)
2. [Features](#features)
3. [Installation](#installation)
4. [Usage](#usage)
5. [Language Basics](#language-basics)
6. [Example](#example)
7. [Roadmap](#roadmap)
8. [Contributing](#contributing)
9. [License](#license)

---

## Overview

Phaeton is a minimalist dynamically-typed language designed for modeling multi-stage conditional logic. Originally created to explore age-based life phase categorization (hence the name – suggesting _phases_ and _iterative refinement_), it now serves as a learning tool for implementing basic interpreter concepts.

---

## Features

- **Stage-based conditionals**: Nested `if/else` chains for multi-tiered logic
- **Dynamic typing**: Automatic type inference for variables
- **Truthy/falsy evaluation**: Built-in rules for conditional checks
- **Simple error reporting**: Clear line-numbered feedback

---

## Installation

```bash
# Clone repository
git clone https://github.com/Stan-breaks/phaeton.git
cd phaeton

# Build interpreter
go build -o phaeton cmd/main.go

# Verify installation
./phaeton --version
```

---

## Usage

```bash
# Run a script
./phaeton run example.phn

# Parse without execution (AST preview)
./phaeton parse example.phn
```

---

## Language Basics

### Variables and Types

```phn
var age = 25;           // Integer
var name = "Alice";      // String
var score = 97.5;     // Float
```

### Conditionals

```phn
if (score >= 90) {
    grade = "A";
} else if (score >= 80) {
    grade = "B";
} else {
    grade = "C";
}
```

### Truthiness Rules

| Value               | Evaluation |
| ------------------- | ---------- |
| `0`                 | Falsy      |
| `""` (empty string) | Falsy      |
| `"nil"`             | Falsy      |
| Non-zero number     | Truthy     |
| Non-empty string    | Truthy     |

---

## Example

**life_stages.phn**

```phn
var stage = "unknown";
var age = 25;

if (age < 18) {
    stage = "minor";
} else if (age < 30) {
    stage = "young adult";
} else if (age < 65) {
    stage = "adult";
} else {
    stage = "senior";
}

print "Life stage: " + stage;
```

Output:

```bash
Life stage: young adult
```

---

## Roadmap

- [x] Basic conditional logic
- [ ] Looping constructs
- [ ] Function definitions
- [ ] Standard library (math, string utils)

---

## Contributing

This project welcomes:

- Bug reports via issues
- Code improvements via pull requests
- Documentation enhancements

Guidelines:

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Submit PR with clear description

---

## License

MIT License - see [LICENSE](LICENSE)
