# go-logfilter

Filter output from servers and perform actions based on matching inputs


## Usage


## Config

```
rules:
  - name: test
    match: test text
    command: echo "%s" >> matches.log
  - name: skip
    match: don't log it
    skip: true
```
