# FAV4

FAV4 is a simple microservice for loading website favicons.

## Usage

The image has almost no configuration, so it's as simple as:

```shell
docker run -p 8080:8080 ghcr.io/snakdy/fav4:main
curl http://localhost:8080/?site=https://github.com
```
