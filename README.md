# Go Ranking
Project created with the objective of learning and put in practice the stuff
learned.

## Dependencies
The dependencies are defined on `go.mod`.
- github.com/stretchr/testify v1.3.0
- github.com/urfave/cli v1.20.0

## Roadmap
- [X] Implement Glicko2.
- [X] Define which type of sport will be the focus of the application of the
  Glicko2. (CS:GO)
- [X] Implement fetcher to get basic Competitors/Matches information.
    - [X] Decide main source of information.
        - [TheScore](https://www.thescoreesports.com/csgo)
- [ ] Try the ranking with real data.
    - [ ] Issue found: Needs to do one extra step if Team doesn't compete during a period.
- [ ] Create a persistence layer to store the Competitors rankings and periods.
- [ ] Tune Glicko2 formulas to accept importance/difficulty of the tournaments 
  to which the Matches belongs.

## Tests
The tests created at this moment reflect the example show on the glicko2 
documentation that can be found on on the root of 
this project ([link](./glicko2.pdf)).

### Observations
- On some tests you may notice that instead of checking for equals, I'm looking
  for an Abs() for a specific margin. That is because the tests are done based
  on the data from the [specification](./glicko2.pdf), which contains multiple
  `floats round`, that doesn't make sense to do on code. So the tests with 
  margins are to consider that small difference.
- Multiple variables referenced on the [specification](./glicko2.pdf) can be
  found on comments on the code with the prefix `doc-ref`.

### Execute ranking cmd
```
cd ./ranking-go
go run ./cmd/ranking.go
```

### Execute the tests
```
cd ./ranking-go
go test ./...
```

## Development
### Stack
- [VSCode](https://code.visualstudio.com/)
    - [Bingo Language Server](https://github.com/saibing/bingo)
  
### VSCode Configurations
Workspace `settings.json`
```json
{
    "go.alternateTools": {
        "go-langserver": "bingo",
      },
      "go.languageServerExperimentalFeatures": {
        "format": true,
        "autoComplete": true
      },
      "go.useLanguageServer": true
}
```

`tasks.json`
```json
{
    "version": "2.0.0",
    "tasks": [
        {
            "label": "Run Tests",
            "type": "shell",
            "options": {
                "cwd": "${workspaceFolder}"
            },
            "command": "go test ./..."
        },
        {
            "label": "Tidy Modules",
            "type": "shell",
            "options": {
                "cwd": "${workspaceFolder}",
                "env": {
                    "GO111MODULE": "on"
                }
            },
            "command": "go mod tidy"
        }
    ]
}
```