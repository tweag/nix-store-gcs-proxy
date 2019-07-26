# nix-store-gcs-proxy - A HTTP nix store that proxies requests to Google Storage

Nix supports multiple store backends such as file, http, s3, ... but not
Google Storage.

Here we provide a http store backend for nix, that will proxy all the reads
and writes to Google Storage.

## Usage

Make sure to have the google credentials installed in `~/.config/gcloud` or
the `GOOGLE_APPLICATION_CREDENTIALS` environment variable.

Start the server in one terminal: `./nix-store-gcs-proxy --bucket-name
<name-of-your-bucket>`

Then in another terminal, use `nix copy --to http://localhost:3000?secret-key=path/to/secret.key <INSTALLABLE>`. Eg:

```sh
$ nix-store --generate-binary-cache-key cache1.example.org cache.key cache.pub
$ nix copy --to http://localhost:3000?secret-key=$PWD/cache.key nixpkgs.hello
```

## TODO

* Section that explains how to setup GCS with the LB CDN.

## License

This work is licensed under the Apache License 2.0.
See [LICENSE](LICENSE) for more details.

## Sponsors

This work has been sponsored by [Digital Asset](https://digitalasset.com) and [Tweag I/O](https://tweag.io).

[![Digital Asset](https://avatars1.githubusercontent.com/u/9829909?s=200&v=4)](http://digitalasset.com)
[![Tweag I/O](https://avatars1.githubusercontent.com/u/6057932?s=200&v=4)](https://tweag.io)

This repository is maintained by [Tweag I/O](http://tweag.io)

Have questions? Need help? Tweet at
[@tweagio](http://twitter.com/tweagio).
