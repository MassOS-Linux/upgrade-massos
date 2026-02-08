This is the `/etc/upgrade-massos/trusted-keys/` directory.

GPG public key (`.asc`) files you want to trust for MassOS builds go here.

By default it contains the default keys used to sign official MassOS builds.
This is suitable if you use the default/official MassOS upgrade server(s).

If you use a custom server for online updates, you may need to place any keys
used to sign builds from that server into this directory. Otherwise validation
of said builds could fail.

Filenames themselves in this directory do not matter, but they should all end
in `.asc` (i.e. they should be ASCII-armored).

If you are the maintainer of a custom MassOS upgrade server, be sure **NEVER**
to place your private signing key(s) in this directory! This directory is for
public keys only. Private keys should always be kept private.
