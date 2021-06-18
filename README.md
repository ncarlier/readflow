# readflow

[![Build Status](https://github.com/ncarlier/readflow/actions/workflows/build.yml/badge.svg)](https://github.com/ncarlier/readflow/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ncarlier/readflow)](https://goreportcard.com/report/github.com/ncarlier/readflow)
[![Docker pulls](https://img.shields.io/docker/pulls/ncarlier/readflow.svg)](https://hub.docker.com/r/ncarlier/readflow/)

Read your Internet article flow in one place with complete peace of mind and freedom.

![Logo](readflow.svg)

## Features

- Read articles from anywhere in one place.
- Save articles for offline reading.
- Create categories and classify new articles automatically thanks to a customizable rule engine.
- External service integration thanks to incoming and outgoing webhooks ([Keeper][keeper], [Pocket][pocket], [Wallabag][wallabag], custom...).
- Receive notifications when new articles are to be read.
- Good user experience on mobile devices.
- No tracking and no publicity.

## Installation

Run the following command:

```bash
$ go get -v github.com/ncarlier/readflow
```

**Or** download the binary regarding your architecture:

```bash
$ curl -sf https://gobinaries.com/ncarlier/readflow | sh
$ # or
$ curl -s https://raw.githubusercontent.com/ncarlier/readflow/master/install.sh | bash
```

**Or** use Docker:

```bash
$ docker run -d --name=readflow ncarlier/readflow
```

## Configuration

Readflow can be configured by using command line parameters or by setting environment variables.

Type `readflow -h` to display all parameters and related environment variables.

All configuration variables are described in [etc/default/readflow.env](./etc/default/readflow.env) file.

## UI

You can access Web UI on http://localhost:8080/ui

![Screenshot](screenshot.png)

## Documentation

The documentation can be found here: https://about.readflow.app/docs

## GraphQL API

You can explore the server API using GraphiQL endpoint: http://localhost:8080/graphiql

## Development

To be able to build the project you will need to:

- Install `makefiles` external helpers:
  ```bash
  $ git submodule init
  $ git submodule update
  ```

Then you can build the project using make:

```bash
$ make
```

Type `make help` to see other possibilities.

## Bakers

These amazing people have sponsored this project:

[![Code Lutin](landing/public/img/code-lutin.svg)](https://www.codelutin.com/)

[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.me/nunux)

***

## License

Readflow is provided under the [GNU Affero General Public License Version 3 (AGPLv3)](https://github.com/ncarlier/readflow/blob/master/LICENSE).

```text
Readflow is a personal news reader service.

Copyright (C) 2021 Nicolas Carlier

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>
```

---

[keeper]: https://keeper.nunux.org
[wallabag]: https://www.wallabag.org
[pocket]: https://getpocket.com/
