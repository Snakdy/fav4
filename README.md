# FAV4

FAV4 is a simple microservice for loading website favicons.

## How it works

The favicon specification is unclear and websites often provide their favicon in a bunch of different ways.
Fav4 attempts to find the favicon by using a collection of "loaders" to attempt to find it.
These loaders are:

1. Attempt to download the `/favicon.png` and `/favicon.ico` files
2. Download and scan the `index.html` file for a `favicon.*` file

## Usage

The image has almost no configuration, so it's as simple as:

```shell
docker run -p 8080:8080 ghcr.io/snakdy/fav4:main
```

The application will accept requests on any path as long as the `site` URL parameter is present:

```shell
curl http://localhost:8080/?site=https://github.com
```
