ssh-stream
=

Allows you to fire off multiple ssh commands and will process the output of those commands.
Allows you to filter the contents of the commands.

Example
===

    ./ssh-stream --filter=cron \
    "me@server1 tail -f /server1/cron/log/file"
    "me@server2 tail -f /server2/cron/log/file"
    "me@server3 tail -f /server3/cron/log/file"

Commands
===

* --filter - Regexp pattern to filter by.
* --user - Username for ssh connection (will prompt if not supplied)
* --password - Password for ssh connection (will prompt if not supplied, silent)
