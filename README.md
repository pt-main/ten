# Ten

A powerful template engine with expressive syntax and full programmatic control. Use it as a standalone CLI tool or integrate it into your Go applications.

## Features

- **Flexible syntax** with three core blocks: config, code, and placeholders
- **Conditional logic** with `if/else`, flags, and variable scoping  
- **String operations** including case conversion, slicing, concatenation
- **CLI-first design** with file processing and documentation generation
- **Go library** for seamless integration via `TemplateReplace` API
- **Customizable** – config blocks allow runtime syntax modification

## Installation

### Quick Start

Clone and run immediately:

```bash
git clone https://github.com/pt-main/ten.git
cd ten
./bin/ten-linux-amd64 process tests/test1.ten name:John vip evening
```

Pre-built binaries for Linux, macOS, and Windows (amd64/arm64) are in `/bin`.

### Build from Source

```bash
cd ten/cmd/cli
go build -o ten .
./ten process tests/test1.ten --name=John --vip --evening
```

Requires Go 1.24.13 or later.

## Usage

### Command Line

#### Process

Transform a template file with flags and placeholders:

```bash
ten process <template_file> [placeholders] [flags] [--O|--OUT=output_file]
```

**Examples:**

```bash
# Print result to stdout
ten process template.ten --name=Alice --prefix=Hello

# Save to file
ten process template.ten --name=Bob --OUT=output.txt

# With flags
ten process template.ten --user=admin --debug --production --OUT=result.txt
```

**Arguments:**
- `template_file` – path to `.ten` template
- `--name=value` – placeholder assignments (lowercase keys)
- `--flag` – boolean flags (lowercase keys, no value)
- `--O=file` or `--OUT=file` – write output to file instead of stdout

#### Docs

Display embedded documentation from a template:

```bash
ten docs template.ten
```

Extracts and prints the `documentation` field from the config block.

### Go API

Use Ten as a library in your Go code:

```go
import (
  "fmt"
  "log"

  "github.com/pt-main/ten/ten"
)

func main() {
  templateString := "var: [[? var ]]"
  result, err := ten.TemplateReplace(
      templateString,
      []string{"flag1", "flag2"},                    // flags
      map[string]string{"var": "value"},             // placeholders
      nil,                                           // config (uses defaults)
  )
  if err != nil {
      log.Fatal(err)
  }
  fmt.Println(result)
}
```

## Template Syntax

A template consists of three block types:

### Config Block

Defines engine behavior and custom documentation:

```
((?CONFIG
    key1 value1
    key2 value2
    documentation Your docs here
))
```

Default keys:
- `pre_plus` – prefix marker (default: `pre+`)
- `post_plus` – suffix marker (default: `post+`)
- `raw` – raw string prefix (default: `raw:`)
- `string_start` / `string_end` – multiline string delimiters (default: `/ss(`/`)/se`)
- `comment` – code comment prefix (default: `//`)
- `ALPHA` – enable nested placeholder support

### Code Block

Execute commands to build variables and logic:

```
{{?CODE
    set varname = value
    logic flagname = is_set
    compare var1 == var2 = result
    place expression = varname
    print text\n
    add var1 var2 = combined
    make upper text = uppercase
}}
```

**Commands:**

| Command | Syntax | Description |
|---------|--------|-------------|
| `set` | `set name = string` | Assign variable. Support multiline string defining. |
| `logic` | `logic flagname = result` | Check if flag exists |
| `compare` | `compare val1 == val2 = result` | Compare strings (`==`/`equals`, `.>`/`startswith`, `>.`/`endswith`) |
| `bool` | `bool var1 and\or var2 = result` | Logical AND/OR |
| `place` | `place expr = name` | Evaluate and store placeholder |
| `add` | `add val1 val2 ... = result` | Concatenate values |
| `make` | `make upper\lower\slice_s:N\slice_e:N val = result` | Transform string |
| `print` | `print text` | Output to stdout |
| `printv` | `printv logic\ts [name]` | Debug print variables from logic scope / template sources scope |
| `//` | Comment | Ignore line |

### Placeholder Block

Insert computed values into output:

```
[[? varname ]]
[[? var if flagname else default_var ]]
[[? var for 3 ]]
```

**Syntax:**
- `[[? variable ]]` – simple substitution
- `[[? var if condition ]]` – conditional (condition: flag, `logic varname`, `has varname`, prefixed with `not`)
- `[[? var if cond else fallback ]]` – with fallback
- `[[? var for N ]]` – repeat N times
- `[[? var1 if cond1 else var2 if cond2 else ... ]]` - nested conditions (if ALPHA mode enable in config)

## Example

**template.ten:**

```
((?CONFIG
    documentation Greeting generator
)){{?CODE
    set greeting = Hello
    make upper user = user_caps
    compare user_caps == nil = isEmpty
    place raw:User if logic isEmpty else user_caps = block1
    place greeting if not has msg else msg = msg
    add msg raw:,\s block1 = result
}}[[? result ]]!
```

**Run:**

```bash
$ ten process template.ten --user:john
Hello, JOHN!

$ ten process template.ten
Hi, User!
```

## Configuration Defaults

| Key | Default | Purpose |
|-----|---------|---------|
| `pre_plus` | `pre+` | Spacing prefix marker |
| `post_plus` | `post+` | Spacing suffix marker |
| `raw` | `raw:` | Raw string prefix (no substitution) |
| `string_start` | `/ss(` | Multiline string start |
| `string_end` | `)/se` | Multiline string end |
| `comment` | `//` | Code comment prefix |
| `ALPHA` | `` | Enable nested placeholder logic |
| `documentation` | `has no docs` | Embedded help text |

## Special Markers

- **`pre+`** – Add space before value
- **`post+`** – Add space after value
- **`raw:`** – Treat string literally (no support spaces, use `\n` and `\s` escapes for newlines and spaces)

## Requirements

- Go 1.24.13 or later (for building from source)
- `lc`, `tap`, `regexp2`

## License

MIT License – see LICENSE file for details

## Contributing

Feedback and contributions welcome. Check issues for open tasks or propose improvements.

---

**Built by Pt** | Template Engine
