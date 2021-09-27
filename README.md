## Build
```
go build
```

## Usage
```
butcher [OPTIONS] address

Application Options:
  -c, --config= config path

Help Options:
  -h, --help    Show this help message
```

## Example
```
./butcher -c config.json 0.0.0.0:53
```

**Docker:**
```
docker run --name butcher --restart always -d -v $(pwd)/config.json:/config.json:ro -p 53:53/udp doorbash/butcher
```
