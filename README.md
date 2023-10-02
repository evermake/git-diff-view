# 🤖 Git-diff server

## Requirements

- Go >= 1.21.1
- Git >= 2.38

## Quick start

Make sure you have Go version at least 1.21. You can check it by running

```shell
go version
# go version go1.21.1 darwin/arm64
```

Clone the repo and start the server (on port `7777` by default with repo in the current working directory)

```shell
git clone https://github.com/evermake/git-diff-view.git
go run ./cmd/app
# Now you can connect your client 🏆
```

## Usage

```
Usage of git-diff-server:
  -port string
    	port to listen on (default "7777")
  -repo string
    	path to the git repository (default is the current directory)
```

## Contributors

Ordered by impact of their contribution:

1. Vadim Z. ([metafates](https://github.com/metafates))
2. Albert A. ([TatarinAlba](https://github.com/TatarinAlba))
3. Vladislav D. ([evermake](https://github.com/evermake))
4. Andrey S. ([5ur3](https://github.com/5ur3))
5. Andrey P. ([majorro](https://github.com/majorro))
