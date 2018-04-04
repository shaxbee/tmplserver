# Template Server

Server for serving static files and templates.  
Templates are parsed using text/template.

## Example

Expose server on port 80 using example templates and data.

```sh
docker run -p 80:80 -v "${PWD}/example:/data"  -e "PAGE_TITLE=Hello World" shaxbee/tmplserver -base /data -env -data data.yaml 
```

Check if it works:

```sh
curl http://localhost
```

Output:

```html
<html>
<head>
    <title>Hello World</title>
</head>
<body>
    <h1>Pilots</h1>
    <ul>
        <li>Leela Turanga</li>
        <li>Bender Bending Rodriguez</li>
    </ul>
</body>
```

## Install

```sh
go get -u github.com/shaxbee/tmplserver/cmd/tmplserver
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

`Env` contains map of environment variables if `-env` is used.  
If value contains comma it is split into slice of strings.

`Data` contains map loaded from yaml or json file specified by `-data`.


