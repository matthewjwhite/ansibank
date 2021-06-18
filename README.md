# ansibank

Ever find yourself wishing you could review the output from a playbook you ran a few weeks ago?

... Well, now you can!

## Description

`ansibank` is a drop-in replacement for `ansible-playbook` that preserves your playbook output.

`ansibank` streams your `ansible-playbook` output as you would expect, while writing it to a local SQLite
database upon completion. AWX/Tower are nifty because they store playbook output for later review -
this project seeks to CLI-ify output preservation in a way that is closer to typical `ansible-playbook` usage.

## Usage

### Build

Simply clone this repository and run `go install .`. Ansibank can then be accessed via `ansibank`.

### Running Playbook

Run your playbook as you would `ansible-playbook`, but replaced with `ansibank`.

For example:
```
ansible-playbook -e SOME_VAR=test my-playbook.yml
```

Would become:
```
ansibank -e SOME_VAR=test my-playbook.yml
```

The only caveat is that the code currently assumes your playbook path is the _last_ argument to `ansibank`.
Note that the playbook path, as with `ansible-playbook`, **does not** need to be an absolute path.

### Listing Previous Runs

Run `ansibank list` from the directory containing your database.

## TODO

* Use TUI for listing, and print selected playbook execution to print output.
* Currently, playbook runs are identified by the playbook path (realpathed). If you move your playbook to
  a different location, it would be helpful to be able to reflect that in the database. The thought is to
  provide a `move` command here.
* The Ansibank DB is created in the directory from where you run `ansibank`, with the name `ansibank-db`.
  Making this configurable would be nice.
