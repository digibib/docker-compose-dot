# docker-compose-dot

Generate Graphviz dot files from docker-compose yaml files.
Inspired by https://github.com/abesto/docker-compose-graphviz
Adapted to use docker-compose yaml v2 and table formatting

## Usage

```
go get github.com/digibib/docker-compose-dot
```

```
Usage:
  docker-compose-dot docker-compose.yml
```

## Docker image use

```
export TAG=21af6b4fd714903cebd3d4658ad35da4d0db0051
```

```
docker pull digibib/docker-compose-dot:$TAG
```

converting a docker-compose.yml in the current dir:

```
docker run --rm -v $(pwd):/tmp digibib/docker-compose-dot:$(TAG) ./app /tmp/docker-compose.yml 2> /dev/null 1> docker-compose.dot
```

You will need the Graphviz package to convert dot to image formats.

#### MIT License

Copyright © 2016 Oslo Public Library <digibib@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
'Software'), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
