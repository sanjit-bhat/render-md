# render-md

A simple wrapper around the [GitHub Markdown API](https://docs.github.com/en/rest/markdown/markdown) that renders a `.md` file to `.html`.

## Prerequisites

- Go
- `gh` CLI, authenticated (`gh auth status`)

## Usage

```sh
go install github.com/sanjit-bhat/render-md@latest
render-md file.md           # writes file.html
open file.html
```

Optional output path:

```sh
render-md file.md -output /tmp/out.html
```
