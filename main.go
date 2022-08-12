package main

import (
	"fmt"
	"io/ioutil"
	"os"

	geojson "github.com/paulmach/go.geojson"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	defaultGridSizeLat float64 = 0.00004291534423 * 3 // geohash9-height x3
	defaultGridSizeLon float64 = 0.00004291534423 * 3 // geohash9-width x3
)

var (
	total   = 0
	failed  = 0
	success = 0

	minlat  = kingpin.Flag("minlat", "min lat").Required().Float64()
	minlon  = kingpin.Flag("minlon", "min lon").Required().Float64()
	maxlat  = kingpin.Flag("maxlat", "max lat").Required().Float64()
	maxlon  = kingpin.Flag("maxlon", "max lon").Required().Float64()
	sizelat = kingpin.Flag("size-lat", "grid size: height(lat)").Default(fmt.Sprintf("%v", defaultGridSizeLat)).Float64()
	sizelon = kingpin.Flag("size-lon", "grid size: width(lon)").Default(fmt.Sprintf("%v", defaultGridSizeLon)).Float64()
	outfile = kingpin.Flag("out", "output file").Required().String()
)

func main() {
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	fc := geojson.NewFeatureCollection()
	id := int64(1)
	for lat := *minlat; lat < *maxlat; lat += *sizelat {
		for lon := *minlon; lon < *maxlon; lon += *sizelon {
			box := [][][]float64{
				{
					{
						lon,
						lat,
					},
					{
						lon,
						lat + *sizelat,
					},
					{
						lon + *sizelon,
						lat + *sizelat,
					},
					{
						lon + *sizelon,
						lat,
					},
					{
						lon,
						lat,
					},
				},
			}

			f := geojson.NewPolygonFeature(box)
			f.SetProperty("id", id)
			fc.AddFeature(f)
			id++
		}
	}

	json, err := fc.MarshalJSON()
	if err != nil {
		fmt.Printf("feature collection marshall json error: %v\n", err)
		os.Exit(1)
	}

	if err := ioutil.WriteFile(*outfile, json, os.ModePerm); err != nil {
		fmt.Printf("write file error: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
