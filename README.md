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

## Usage

Simply clone this repository and run `go install .`.

From here, run your playbook as you would with with `ansible-playbook`, but with `ansible-playbook`
replaced with `ansibank`.

## TODO

* Running queries against your SQLite DB directly isn't too helpful. The plan is to add commands to list
  your playbook runs and pick a particular one to view. Hopefully with a nifty TUI.
* Currently, playbook runs are identified by the playbook path (realpathed). If you move your playbook to
  a different location, it would be helpful to be able to reflect that in the database. The thought is to
  provide a `move` command here.
