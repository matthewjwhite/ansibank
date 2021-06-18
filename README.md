# ansibank

Ever find yourself wishing you could review the output from a playbook you ran a few weeks ago?

... Well, now you can!

## What?

`ansibank` is a drop-in replacement for `ansible-playbook`.

`ansibank` streams your `ansible-playbook` output as you would expect, while writing it to a local SQLite
database upon completion.

## Why?

Having some experience with Tower, I particularly :heart:ed the ability to preserve playbook logs. When
developing playbooks, I generally run them multiple times to debug, sometimes across different windows,
sometimes closing those windows between runs. It's often helpful to review those logs, and losing them
can be a huge pain.

`ansibank` seeks to totally eliminate the possibility of losing those logs.

# Usage

## Build

Simply clone this repository and run `go install .`. Ansibank can then be accessed via `ansibank`.

## Running Playbook

Run your playbook as you would `ansible-playbook`, but replaced with `ansibank`.

For example:
```
ansible-playbook -e SOME_VAR=test my-playbook.yml
```

Would become:
```
ansibank -e SOME_VAR=test my-playbook.yml
```

The only caveat is that the code currently assumues your playbook path is the _last_ argument to `ansibank`.
Note that the playbook path, like with `ansible-playbook`, **does not** need to be an absolute path.

## Listing Previous Runs

Run `ansibank list` from the directory containing your database.

# TODO

* Use TUI for listing, and print selected playbook execution to print output.
* Currently, playbook runs are identified by the playbook path (realpathed). If you move your playbook to
  a different location, it would be helpful to be able to reflect that in the database. The thought is to
  provide a `move` command here.
* The Ansibank DB is created in the directory from where you run `ansibank`, with the name `ansibank-db`.
  Making this configurable would be nice.
