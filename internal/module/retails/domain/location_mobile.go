package domain

import (
	"github.com/kamva/mgm/v3"
	"samm/pkg/utils"
	"time"
)

type Status struct {
	SnoozeTo string `json:"snooze_to"`
	Status   string `json:"status" bson:"status"`
	Meta     Meta   `json:"meta" bson:"meta"`
}
type Meta struct {
	NameEN string `json:"name_en" bson:"name_en"`
	NameAr string `json:"name_ar" bson:"name_ar"`
	Color  string `json:"color" bson:"color"`
}

type LocationMobile struct {
	mgm.DefaultModel `bson:",inline"`
	Name             Name          `json:"name" bson:"name"`
	City             City          `json:"city" bson:"city"`
	Street           Name          `json:"street" bson:"street"`
	CoverImage       string        `json:"cover_image" bson:"cover_image"`
	Logo             string        `json:"logo" bson:"logo"`
	SnoozeTo         *time.Time    `json:"snooze_to" bson:"snooze_to"`
	IsOpen           bool          `json:"is_open" bson:"is_open"`
	WorkingHour      []WorkingHour `json:"-" bson:"working_hour"`
	Phone            string        `json:"phone" bson:"phone"`
	Coordinate       Coordinate    `json:"coordinate" bson:"coordinate"`
	BrandDetails     BrandDetails  `json:"brand_details" bson:"brand_details"`
	PreparationTime  int           `json:"preparation_time" bson:"preparation_time"`
	Distance         float64       `json:"distance" bson:"distance"`
	Country          Country       `json:"country" bson:"country"`
	Status           Status        `json:"status" bson:"-"`
	//NextEventTime    string        `json:"next_event_time" bson:"-"`
}

func (payload *LocationMobile) SetOpenStatus() {
	open := payload.IsOpen
	now := time.Now().UTC()

	if open {
		payload.Status = Status{
			SnoozeTo: "",
			Status:   "open",
			Meta: Meta{
				NameEN: "Open",
				NameAr: "مفتوح",
				Color:  "",
			},
		}
	} else {
		payload.Status = Status{
			SnoozeTo: "",
			Status:   "closed",
			Meta: Meta{
				NameEN: "Closed",
				NameAr: "مغلق",
				Color:  "",
			},
		}
	}

	if payload.SnoozeTo != nil && now.Before(*payload.SnoozeTo) {
		payload.Status = Status{
			SnoozeTo: "",
			Status:   "busy",
			Meta: Meta{
				NameEN: "Busy",
				NameAr: "مشغول",
				Color:  "",
			},
		}
	}

}
func (payload *LocationMobile) SetDistance(lat float64, lng float64) {
	if lat != 0 && lng != 0 {
		payload.Distance = utils.Distance(payload.Coordinate.Coordinates[1], payload.Coordinate.Coordinates[0], lat, lng)
	}
}

//func (payload *LocationMobile) SetNextEvent() {
//	daysOfWeek := map[string]int{
//		"sunday":    0,
//		"monday":    1,
//		"tuesday":   2,
//		"wednesday": 3,
//		"thursday":  4,
//		"friday":    5,
//		"saturday":  6,
//	}
//	//targetLocation, _ := time.LoadLocation(payload.Country.Timezone)
//
//	currentTime := time.Now().UTC()
//	currentDay := strings.ToLower(currentTime.Weekday().String())
//	currentHour := currentTime.Format("15:04:05")
//	workingHBD := WorkingHourByDay(payload.WorkingHour)
//
//	if payload.IsOpen {
//		// Get The Current Period
//		currentPeriod := GetCurrentPeriod(workingHBD[currentDay], currentTime)
//		// if full day return null for open 24 hours
//		if currentPeriod != nil && currentPeriod.IsFullDay {
//			payload.NextEventTime = ""
//			return
//		}
//
//		//closeTime, _ := time.Parse("15:04:05", currentPeriod.To)
//		//
//		//eventCloseTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), closeTime.Hour(), closeTime.Minute(), 0, 0, currentTime.Location())
//		//
//		//if currentPeriod.To {
//		//
//		//}
//		// check the end if it near to end of day then check the next day to handle the close in the next day
//
//	} else {
//
//	}
//
//	var nextEventTime time.Time
//	var nextEventDescription string
//
//	for _, wh := range payload.WorkingHour {
//		openTime, _ := time.Parse("15:04:05", wh.From)
//		closeTime, _ := time.Parse("15:04:05", wh.To)
//
//		whDay, ok := daysOfWeek[wh.Day]
//		if !ok {
//			fmt.Println("invalid day: ", wh.Day)
//			return
//		}
//
//		eventOpenTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), openTime.Hour(), openTime.Minute(), 0, 0, currentTime.Location())
//		eventCloseTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), closeTime.Hour(), closeTime.Minute(), 0, 0, currentTime.Location())
//
//		if openTime.After(closeTime) {
//			// Handles transition to the next day for closing time
//			eventCloseTime = eventCloseTime.Add(24 * time.Hour)
//		}
//
//		if wh.IsFullDay {
//			if daysOfWeek[currentDay] < whDay || (daysOfWeek[currentDay] == whDay && currentHour <= "00:00") {
//				if nextEventTime.IsZero() || eventOpenTime.Before(nextEventTime) {
//					nextEventTime = eventOpenTime
//					nextEventDescription = fmt.Sprintf("Next open time on %s at 00:00", wh.Day)
//				}
//			}
//		} else {
//			if daysOfWeek[currentDay] < whDay || (daysOfWeek[currentDay] == whDay && currentHour <= wh.From) {
//				if nextEventTime.IsZero() || eventOpenTime.Before(nextEventTime) {
//					nextEventTime = eventOpenTime
//					nextEventDescription = fmt.Sprintf("Next open time on %s at %s", wh.Day, wh.From)
//				}
//			}
//
//			if daysOfWeek[currentDay] < whDay || (daysOfWeek[currentDay] == whDay && currentHour <= wh.To) {
//				if nextEventTime.IsZero() || eventCloseTime.Before(nextEventTime) {
//					nextEventTime = eventCloseTime
//					nextEventDescription = fmt.Sprintf("Next close time on %s at %s", wh.Day, wh.To)
//				}
//			}
//		}
//	}
//
//	if nextEventDescription == "" {
//		return
//	}
//	payload.NextEventTime = nextEventTime.Format("15:04:05")
//
//}
//
//func WorkingHourByDay(workingHours []WorkingHour) map[string][]WorkingHour {
//	result := make(map[string][]WorkingHour, 0)
//	var query []Group
//	From(workingHours).GroupByT(
//		func(i WorkingHour) string {
//			return i.Day
//		},
//		func(i WorkingHour) WorkingHour {
//			return i
//		}).ToSlice(&query)
//
//	for _, resultItem := range query {
//		workingDays := make([]WorkingHour, 0)
//		for _, workingDay := range resultItem.Group {
//			workingDays = append(workingDays, workingDay.(WorkingHour))
//		}
//		result[resultItem.Key.(string)] = workingDays
//	}
//	return result
//}
//
//func GetCurrentPeriod(workingHours []WorkingHour, currentTime time.Time) *WorkingHour {
//
//	for _, item := range workingHours {
//		if item.IsFullDay {
//			return &item
//		}
//		openTime, _ := time.Parse("15:04:05", item.From)
//		closeTime, _ := time.Parse("15:04:05", item.To)
//
//		eventOpenTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), openTime.Hour(), openTime.Minute(), 0, 0, currentTime.Location())
//		eventCloseTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), closeTime.Hour(), closeTime.Minute(), 0, 0, currentTime.Location())
//		if eventOpenTime.Before(eventCloseTime) && currentTime.After(eventOpenTime) && currentTime.Before(eventCloseTime) {
//			return &item
//
//		}
//		if eventOpenTime.After(eventCloseTime) && currentTime.Before(eventOpenTime) && currentTime.After(eventCloseTime) {
//			return &item
//		}
//	}
//	return nil
//}
