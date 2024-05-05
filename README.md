# binhost

A Gentoo binary host platform. Enables building and tracking various
targets.

## Table of Contents

<!-- toc -->
- [API](#api)
  - [<code>POST /v1/upload</code>](#post-v1upload)
  - [<code>GET /v1/targets</code>](#get-v1targets)
  - [<code>POST /v1/targets/:target</code>](#post-v1targetstarget)
- [License](#license)
<!-- /toc -->

## Developing

We use [mise]() to manage development versions and tasks.

```bash
mise run build
```

Now you can run the server with:

```bash
ENV=development ./bin/binhost
```

## API

Loose documentation of the API provided by `binhost` is below.

### `POST /v1/upload`

Uploads the provided `gpkg`. If a `target` is not specified via URL
parameters, the `CHOST` is used as the target name. Errors if a target
doesn't exist.

### `GET /v1/targets`

Lists all of the available targets (package indexes).

### `POST /v1/targets/:target`

Creates the provided target.

## License

AGPL-3.0
