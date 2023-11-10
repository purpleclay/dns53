---
icon: material/file-code-outline
---

# Compiling from Source

Download both [Go 1.21+](https://go.dev/doc/install) and [go-task](https://taskfile.dev/#/installation). Then clone the code from GitHub:

```sh
git clone https://github.com/purpleclay/dns53.git
cd dns53
```

Build:

```sh
task
```

And check that everything works:

```sh
./dns53 version
```
