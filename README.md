# JavaScript to Go Transpiler

Transpiles JavaScript into Go for faster performance and native binaries.

## Usage

```bash
go run .
```

Reads `case1.js` and outputs a runnable `output.go`.

## Supported Features

- `let` / `var` variable declarations
- Reassignment of variables (typed as `any`)
- `console.log()` → `fmt.Println()`
- Strings and integers

## Example

```js
let a = 10;
a = "test";
console.log(a);
```

Outputs:

```go
package main

import "fmt"

func main() {
    var a any = 10
    a = "test"
    fmt.Println(a)
}
```

## Pipeline

Lexer → Parser → Semantic Analyzer → Code Generator
