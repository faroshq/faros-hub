# Dev IDP

For dev we use Dex as identity provider. More of it: https://github.com/dexidp/dex.git

If you setup kind cluster with `make setup-kind` it will install Dex and configure it to use static users.

You need to make sure `dex.dev.faros.sh` resolved to 127.0.0.1 in your machine.
Via `/etc/hosts` or other means.


Setup `.env` file as per `env.example` file.

Once this is done you can run faros with IDP inside the kind cluster.
