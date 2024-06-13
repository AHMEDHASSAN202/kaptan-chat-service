package utils

type Point struct {
	Lat float64
	Lng float64
}

type Polygon []Point

func EGPolygon() Polygon {
	return Polygon{
		{
			Lat: 31.6720835,
			Lng: 25.0268555,
		},
		{
			Lat: 31.0717559,
			Lng: 24.6972656,
		},
		{
			Lat: 30.5244133,
			Lng: 24.6972656,
		},
		{
			Lat: 29.7834495,
			Lng: 24.4995117,
		},
		{
			Lat: 29.5161104,
			Lng: 24.5874023,
		},
		{
			Lat: 28.6520306,
			Lng: 24.8510742,
		},
		{
			Lat: 27.1960144,
			Lng: 24.8071289,
		},
		{
			Lat: 26.0172976,
			Lng: 24.7412109,
		},
		{
			Lat: 24.9462191,
			Lng: 24.7412109,
		},
		{
			Lat: 23.6042618,
			Lng: 24.8510742,
		},
		{
			Lat: 22.5328537,
			Lng: 24.7631836,
		},
		{
			Lat: 21.6778483,
			Lng: 24.7631836,
		},
		{
			Lat: 21.6982655,
			Lng: 25.3125,
		},
		{
			Lat: 21.6574282,
			Lng: 27.3339844,
		},
		{
			Lat: 21.8207079,
			Lng: 29.1357422,
		},
		{
			Lat: 21.8003081,
			Lng: 30.4541016,
		},
		{
			Lat: 21.7799053,
			Lng: 32.2338867,
		},
		{
			Lat: 21.7390912,
			Lng: 33.6621094,
		},
		{
			Lat: 21.4939636,
			Lng: 34.5849609,
		},
		{
			Lat: 21.4530686,
			Lng: 35.5078125,
		},
		{
			Lat: 21.6574282,
			Lng: 37.0019531,
		},
		{
			Lat: 22.3500758,
			Lng: 37.265625,
		},
		{
			Lat: 27.8390761,
			Lng: 34.2333984,
		},
		{
			Lat: 28.9216313,
			Lng: 34.7607422,
		},
		{
			Lat: 30.4865508,
			Lng: 34.7167969,
		},
		{
			Lat: 31.203405,
			Lng: 34.6289063,
		},
		{
			Lat: 31.4661537,
			Lng: 33.75,
		},
		{
			Lat: 31.5036293,
			Lng: 33.1347656,
		},
		{
			Lat: 31.5410899,
			Lng: 31.8164063,
		},
		{
			Lat: 31.6159659,
			Lng: 30.9814453,
		},
		{
			Lat: 31.5410899,
			Lng: 29.8828125,
		},
		{
			Lat: 31.3161014,
			Lng: 29.0478516,
		},
		{
			Lat: 31.5410899,
			Lng: 27.2900391,
		},
		{
			Lat: 31.6159659,
			Lng: 25.9277344,
		},
		{
			Lat: 31.7281671,
			Lng: 25.3564453,
		},
		{
			Lat: 31.6720835,
			Lng: 25.0268555,
		},
	}
}

func IsInsidePolygon(polygon Polygon, lat float64, lng float64) bool {

	point := Point{
		Lat: lat,
		Lng: lng,
	}
	// Perform ray casting algorithm
	intersections := 0

	for i := 0; i < len(polygon); i++ {
		j := (i + 1) % len(polygon)

		if (polygon[i].Lat > point.Lat) != (polygon[j].Lat > point.Lat) &&
			point.Lng < (polygon[j].Lng-polygon[i].Lng)*(point.Lat-polygon[i].Lat)/(polygon[j].Lat-polygon[i].Lat)+polygon[i].Lng {
			intersections++
		}
	}

	return intersections%2 == 1
}
func GetCountryFromLatLng(lat float64, lng float64) string {
	if IsInsidePolygon(EGPolygon(), lat, lng) {
		return "EG"
	}
	return "SA"
}
