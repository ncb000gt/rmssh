ssh-stream
=

Allows you to fire off multiple ssh commands and will process the output of those commands.
Allows you to filter the contents of the commands.

Example
===

    ./ssh-stream --filter=cron \
    "me@server tail -f /server/cron/log/file"

Commands
===

* --filter - Regexp pattern to filter by.
* --user - Username for ssh connection (will prompt if not supplied)
* --password - Password for ssh connection (will prompt if not supplied, silent)
