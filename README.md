# go-sfomuseum-airports

Go package for working with airports, in a SFO Museum context. 

## Documentation

Documentation is incomplete.

## Tools

### lookup

```
$> ./bin/lookup -h
Usage of ./bin/lookup:
  -lookup-uri string
    	Valid sources are: sfomuseum:// (default "sfomuseum://")
```

Lookup an airport by its IATA or ICAO code.

```
./bin/lookup EGLL
$> LHR EGLL "London Heathrow Airport" 102556703

$> ./bin/lookup YUL
YUL CYUL "Montreal-Pierre Elliott Trudeau International Airport" 102554351
```

## See also

* https://github.com/sfomuseum-data/sfomuseum-data-whosonfirst
