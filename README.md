# Netlify Git API

Turn your Git repository into a REST API.

## Installation

Grab the binary from your platform from the [latest release](https://github.com/netlify/netlify-git-api/releases) and place it somewhere in your PATH.

## Usage

```bash
cd my-project
netlify-git-api users add
netlify-git-api serve
```

This will add a new user and start serving an API for your repo.

## Options

See `netlify-git-api help` for options and sub commands.

## Using with netlify CMS

Make sure to configure the `backend` in your `config.yml` like this:

```yaml
backend:
  name: netlify-api
  url: localhost:8080
```

## License

The MIT License (MIT)

Copyright (c) 2015 MakerLoop Inc

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
