# goyokiki

A simple, prompt-based CLI assistant for periodic notes.

## Overview

`goyokiki` reads a YAML configuration file, presents each question as an interactive prompt, and outputs the completed answers as a Markdown document.

## Usage

```bash
go run ./cmd/goyokiki [config.yaml]
```

If no config file is specified, `config.yaml` in the current directory is used.

You can redirect the Markdown output to a file:

```bash
go run ./cmd/goyokiki example/config.yaml > output.md
```

## Config File Format

Define questions in a YAML file:

```yaml
questions:
  - title: "今日やったこと"
    type: free

  - title: "気分はどうですか"
    type: single
    options:
      - よい
      - ふつう
      - わるい
```

**Question types:**

| Type     | Description                           |
|----------|---------------------------------------|
| `free`   | Free-text input                       |
| `single` | Single-choice selection from a list   |

## Markdown Output Format

```markdown
# 今日やったこと: コードを書いた

# 気分はどうですか: よい
- [x] よい
- [ ] ふつう
- [ ] わるい
```

## Build

```bash
go build -o goyokiki ./cmd/goyokiki
```
