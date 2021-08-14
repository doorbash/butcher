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
Build:
```
docker build -t doorbash/butcher .
```
Run:
```
docker run --restart always -d --name butcher -v $(pwd)/config.json:/config.json:ro -p 53:53/udp doorbash/butcher
```