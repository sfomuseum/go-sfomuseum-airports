# go-sfomuseum-airports

 Go package for working with airports, in a SFO Museum context. 

## Install

You will need to have both `Go` (specifically version [1.12](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make tools
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

### lookup

Lookup an airport by its IATA or ICAO code.

```
./bin/lookup EGLL
$> LHR EGLL "London Heathrow Airport" 102556703

$> ./bin/lookup YUL
YUL CYUL "Montreal-Pierre Elliott Trudeau International Airport" 102554351
```

## See also

* https://github.com/sfomuseum/go-sfomuseum-airports-tools
* https://github.com/sfomuseum-data/sfomuseum-data-whosonfirst
