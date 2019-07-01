# Echo path web server #

A simple localhost golang standard web server at the requested port or with a fallback 8770 port responding with the requested path. 

## Launch: ##

`λ go run echopathws.go 8770`

Starting localhost:8770...

Set HTTP handler @ localhost:8770...

## Example: ##

_Http Get Request:_

`http://localhost:8770/simple web request/`

_Html Response:_

`"simple web request/"`