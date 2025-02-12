# Phaeton Language Interpreter

_Phaeton is a minimalist, dynamically-typed interpreted language designed for modeling stage‐based logic and conditional workflows._

> **Note:** The name “Phaeton” reflects its original intent to model life phases through conditional logic. Over time, it has evolved into a learning tool for building interpreters in Go.

---

## Table of Contents

1. [Overview](#overview)
2. [Features](#features)
3. [Installation](#installation)
4. [Usage](#usage)
5. [Language Basics](#language-basics)
6. [Project Structure](#project-structure)
7. [Detailed Package Documentation](#detailed-package-documentation)
   - [Interpreter Package](#interpreter-package)
   - [Environment Package](#environment-package)
   - [Models Package](#models-package)
   - [Parse Package](#parse-package)
   - [Tokenize Package](#tokenize-package)
   - [Utils Package](#utils-package)
8. [Example](#example)
9. [Roadmap](#roadmap)
10. [Contributing](#contributing)
11. [License](#license)

---

## Overview

Phaeton is designed as a simple interpreted language that supports:

- **Stage-based conditionals:** Nested `if/else` and `else if` constructs for multi-tiered logic.
- **Dynamic typing:** Variables are dynamically typed with automatic type inference.
- **Scoped environments:** Functions and loops maintain their own scopes.
- **Simple error reporting:** Descriptive messages that help pinpoint syntax or runtime issues.

The project is implemented in Go and is organized into several packages, each handling a specific part of the interpreter:

- **`interpreter`** – Orchestrates program execution.
- **`environment`** – Manages variables and function scopes.
- **`models`** – Defines the core data structures (tokens, AST nodes, etc.).
- **`parse`** – Converts tokens into an abstract syntax tree (AST).
- **`tokenize`** – Splits the raw source code into tokens.
- **`utils`** – Contains helper functions for token analysis and other utilities.

---

## Features

- **Conditional Statements:** Supports `if`, `else if`, and `else` blocks.
- **Looping Constructs:** Includes `while` and `for` loops.
- **Functions:** Function definitions and function calls with local scoping.
- **Print Statements:** Output values to the console.
- **Error Handling:** Clear, descriptive error messages for debugging.
- **Dynamic Evaluation:** Expressions are parsed, built into an AST, and evaluated at runtime.

---

## Installation

Clone the repository and build the project using Go:

### Clone the repository

```bash
git clone https://github.com/Stan-breaks/phaeton.git
cd phaeton
```

### Build the interpreter executable

```bash
go build -o phaeton cmd/main.go
```

### Build the interpreter executable

```bash

go build -o phaeton cmd/main.go
```

## Usage

After building the project, you can use the interpreter in various ways:

### Run a script:

```bash

./phaeton run example.phn
```

### Parse without execution (view the AST):

```bash

./phaeton parse example.phn
```

### Tokenize the source (view the tokens):

```bash

./phaeton tokenize example.phn
```

## Language Basics

- Phaeton uses a simple, intuitive syntax. Here are some examples:
  Variables and Types

```phn
var age = 25; // Integer
var name = "Alice"; // String
var score = 97.5; // Float

```

- Conditionals

```phn
if (score >= 90) {
grade = "A";
} else if (score >= 80) {
grade = "B";
} else {
grade = "C";
}
```

- Loops

```phn
// While loop
var count = 0;
while (count < 5) {
print count;
count = count + 1;
}

// For loop (with initialization, condition, and expression)
for (var i = 0; i < 5; i = i + 1) {
print i;
}
```

# Project Structure

The project is organized into the following directories and packages:

- cmd/
  Contains the main entry point (main.go) which sets up command-line options (run, parse, tokenize).

- interpreter/
  Houses the core interpreter logic. It processes tokens, handles control structures (if, while, for), function calls, variable assignments, and more.

- environment/
  Manages variable and function scopes. This package provides methods to push/pop scopes and set, get, or reset variable values.

- models/
  Defines essential data structures such as tokens (TokenInfo), AST nodes (e.g., NumberNode, StringNode, NilNode), and helper structs for statement positions (e.g., IfStatementPositions, ForStatementPositions).

- parse/
  Implements the parsing logic. It converts a stream of tokens into an abstract syntax tree (AST) that the interpreter can evaluate.

- tokenize/
  Converts raw source code into a series of tokens (TokenInfo), each annotated with its type and content.

- utils/
  Provides utility functions for token analysis, such as finding semicolons, matching parentheses/braces, and handling argument lists.

## Detailed Package Documentation

### Interpreter Package

- The interpreter package is the core executor for Phaeton. It includes:

`Interprete(tokens []models.TokenInfo) (interface{}, error)`

- The main entry point that loops over tokens and dispatches control to specific handler functions based on token types.

  - Control Structure Handlers:
  - Functions like handleIf, handleWhile, and handleFor parse and execute conditional and loop constructs by:

    1. Identifying token boundaries.
    2. Parsing expressions.
    3. Managing scope using environment.Global.PushScope() and PopScope().

  - Function Handling:

    1.  Functions such as handleFun and handleFunCall manage function definition and invocation, including parameter passing and local scope management.

  - Expression and Variable Handling:
    1.  handleAssignment and handleReassignment for variable declarations and updates.
    2.  handleExpression parses arithmetic and logical expressions by integrating with the parse package.

Refer to the inline comments in the source for more details on each function’s behavior.

### Environment Package

The environment package is responsible for maintaining the state of variables and functions. It typically exposes methods like:

    - PushScope() / PopScope()
    To manage nested scopes (e.g., within functions or loops).

    - Set(variable, value)
    To define or update a variable in the current scope.

    - Get(variable)
    To retrieve the value of a variable.

    - Reset(variable, value)
    To update an existing variable’s value.

This package ensures that variables in inner scopes do not conflict with global variables.

### Models Package

The models package defines the core data structures used across the interpreter, including:

    - AST Nodes:
    Nodes such as NumberNode, StringNode, and NilNode represent evaluated expressions.

    - TokenInfo:
    Structures that store token details (type, lexeme, position).

    - Statement Positions:
    Helper structs (like IfStatementPositions, ForStatementPositions, etc.) that store the boundaries of different language constructs within the token stream.

### Parse Package

The parse package converts tokens into an abstract syntax tree (AST). Key aspects include:

    - Parsing Expressions:
    Converting arithmetic, logical, or function call expressions into nodes.
    - Error Handling:
    Reporting syntax errors with descriptive messages.
    - AST Construction:
    Producing nodes that are later evaluated by the interpreter.

### Tokenize Package

The tokenize package splits raw source code into a sequence of tokens. It is responsible for:

    - Lexical Analysis:
    Recognizing keywords, identifiers, operators, literals, and punctuation.
    - Error Reporting:
    Indicating any lexical errors (e.g., unrecognized characters).

### Utils Package

The utils package contains helper functions that support the other packages:

    - Token Analysis:
    Functions to find semicolons, matching parentheses or braces, and argument boundaries.
    - General Helpers:
    Utilities to check for specific patterns in tokens (e.g., whether a token represents a function call).

## Example

- Create a file called example.phn with the following content:

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

- Run the interpreter:

```phn

./phaeton run example.phn
```

- Expected Output:

```phn

Life stage: young adult
```

# Roadmap

- [ ] Basic conditional logic
- [ ] Looping constructs (while and for loops)
- [ ] Function definitions and standard library enhancements (math, string utilities)
- [ ] Improved error diagnostics and debugging output

# Contributing

- Contributions to Phaeton are welcome. To contribute:

  1. Fork the repository.
  2. Create a feature branch.
  3. Test your changes with the test provided
  4. Submit a pull request with a detailed description of your changes.

- For bug reports or feature requests, please open an issue on GitHub.

# License

- Phaeton is released under the MIT License.

# Final Notes

- Phaeton serves as a platform for learning about interpreter design, dynamic typing, and scope management in Go. Whether you’re interested in the language itself or the underlying implementation techniques, this project offers a modular and accessible codebase to explore.
- Happy coding!
