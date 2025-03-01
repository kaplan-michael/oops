## Opps default backend

Simple go & chi based default backend for handling errors and not founds. 

### Usage
Default template and errors definitions included is included, but you can mount your own.
> [!NOTE] see [template.tmpl](template.tmpl) for the html template and 
> [errors.yaml](errors.yaml) for the error definitions.

By default, these should be mounted to the root of the container. on the following paths:
`template.tmlp mount to /template.tmpl`
`errors.yaml mount to /errors.yaml`

You can also override the location of the template and errors file by setting the following environment variables:
```shell
TEMPLATE=template.tmpl
ERRORS=errors.yaml
```
You can also set a custom bind port and log level(debug/info) by setting the following environment variables:
```shell
LOGLEVEL=info
PORT=8080
```

Images are prebuilt and available on `quay.io/mkaplan/opps:<version>` or `quay.io/mkaplan/opps:latest`

### Status
somewhat stable.

### License
Apache License 2.0

