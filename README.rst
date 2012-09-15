walseq
======

walseq will remotely work like seq (GNU) or jot (BSD) and will print sequential
Postgresql WAL segment filenames. This can come in handy if you fail at life
and have to operate on a range of WAL segments (copy, compress, transfer,
etc.). Whatever you're doing with this is probably wrong, good luck.

Example::

    $ walseq 0000000100000454000000A1 000000010000045600000014
    0000000100000454000000A1
    0000000100000454000000A2
    [...]
    000000010000045600000013
    000000010000045600000014

Installation::

    go build
    sudo cp walseq /usr/local/bin

