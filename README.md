# Template Server

Server for serving static files and templates.

Templates are parsed using text/template.

## Install

```sh
go get -u github.com/shaxbee/tmplserver/cmd/tmplserver
```

## Example

```sh
PAGE_TITLE=tmplserver tmplserver -base ../example -env -data ../example/data.yaml
```

## Usage

### Arguments:
* `-addr` (`0.0.0.0:8080`) - Listen address
* `-base` (`.`) - Files base path
* `-cert` - TLS cerficate path
* `-key` - TLS key path
* `-env` (`false`) - Load environment variables
* `-data` - Template data file path (json/yaml) 

### Template values

Env contains map of environment variables if `-env` is used.  
If value contains comma it is split into slice of strings.

Data contains map loaded from yaml or json file specified by `-data`.


