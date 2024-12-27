### Sample real-time hardware monitor app with Golang

Sample  autoupdating hardware monitor application in Golang to display data about what is happening on the device.

The app can show the data as cli/cmd app (default) or in web browser via websockets with the help of [HTMX](https://htmx.org).

Data is updated every 5 seconds.

#### Prerequisites
- Golang (used version is 1.23.4)
- Web Browser (Optional) - to see the web interface.

#### Usage

- cli/cmd (default): `go run .\cmd\main.go` or `go run .\cmd\main.go cmd`

- web
    1. `go run .\cmd\main.go ws`
    2. open [localhost](http://localhost:9000) (may need to wait 5 seconds to pull the data from the server)