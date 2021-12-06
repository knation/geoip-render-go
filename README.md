# geoip-render-go
Golang project for looking up geo from an IP. Meant to be deployed on Render.

## Dependencies

In order to use this project, you'll need a copy of your own [Maxmind GeoIP database](https://www.maxmind.com/en/geoip2-services-and-databases). You can sign up for the GeoLite2 database [here](https://www.maxmind.com/en/geolite2/signup?lang=en).

This project then relies on the [oschwald/geoip2-golang](https://pkg.go.dev/github.com/oschwald/geoip2-golang) package for looking up an IP in a MaxMind GeoIP database. See the documentation for all of the queries that can be made and the resulting structs.

## Environment variables

This project relies on the following environment variables:

| ENV Variable | Description                                                                | Required | Default   |
|--------------|----------------------------------------------------------------------------|----------|-----------|
| `GEO_FILE`   | The location of your Maxmind GeoIP database (e.g., `./GeoLite2-City.mmdb`) | Yes      | None      |
| `MODE`       | The mode to launch the application in.                                     | No      | "release" |
| `PORT`       | The port for the web service to listen on.                                 | No      | 3000      |

## How to use

### 1. Fork Repository

Fork this repository into your own account.

### 2. Include your Maxmind GeoIP database

Include your Maxmind GeoIP database in the root of the project. For example, `./GeoLite2-City.mmdb`. Commit this to the repository.

### 3. Update code

Update the `handler()` function in `main.go` to make the appropriate query and return the response.

The example in `main.go` uses the IP address provided to make a `City()` query and returns the zip code as JSON.

### 4. Deploy to render

Create a new web service on Render. Use your new repository and set the above environment variables to make them available at run time.

### 5. Make request

The web service listens for requests to `/:ipaddress`.

## Notes

- This project DOES NOT log web requests. It only prints status logs. If you want to log each request, you can add it to the Go Gin logger.
