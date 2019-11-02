# trello-tui

### A terminal ui for trello

# Usage
```
     -board string
           board name
     -log
           Log to file
     -refresh duration
           refresh interval (min=1s) (default 10s)
     -vv
           Increase verbosity level
```

The following environemnt variables are required to be configured:
```
export TRELLO_USER=user
export TRELLO_KEY=key
export TRELLO_TOKEN=token
```

Run with the following command:
```
go run main.go -refresh=30s -board="Board Name"
```
