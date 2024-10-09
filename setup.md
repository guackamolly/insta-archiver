# setup

To properly run this service, you will need to a setup a `.env` file. Start by creating a copy of the `.env.tpl` file and fill the variables with values appropriate for the execution context.

Then, all you need to do is to run the service with the following command:

```bash
go run cmd/web/web.go
```

## Docker

To build the image:

```bash
./tools/build-docker
```

To run the container (include -d if running as daemon):

```bash
./tools/run-docker
```

## Virtual Host

If you're deploying the service as a virtual host, first run the script to update the webpages relative urls.

```bash
./tool/update-virtual-host
```