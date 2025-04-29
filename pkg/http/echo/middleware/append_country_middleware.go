package echomiddleware

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"kaptan/pkg/utils"
	"strconv"
)

func AppendCountryMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		Lat := req.Header.Get("Lat")
		Lng := req.Header.Get("Lng")

		if Lat != "" && Lng != "" {
			latPoint, errLat := strconv.ParseFloat(Lat, 64)
			lngPoint, errLng := strconv.ParseFloat(Lng, 64)
			if errLat == nil && errLng == nil {
				country := utils.GetCountryFromLatLng(latPoint, lngPoint)
				fmt.Println("Country => ", country)
				c.Request().Header.Set("Country-Id", country)
			}
		}

		return next(c)
	}
}
