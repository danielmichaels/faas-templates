package handlers

import (
	"database/sql"
	"fmt"
	"github.com/rickar/cal/v2"
	"github.com/rickar/cal/v2/au"
	"time"
)

type Time struct {
	Start         string
	End           string
	Income        float64
	TotalTime     int
	TotalTimeCalc time.Duration
}
type TimeModel struct {
	Db *sql.DB
}

type Db struct {
	*sql.DB
}

func NewDB(db *sql.DB) *Db {
	return &Db{db}
}

func NewTimeModel(db *sql.DB) *TimeModel {
	return &TimeModel{Db: db}
}

func (t *Time) InHours() (time.Duration, error) {
	d1, err := time.Parse("2006-01-02 15:04", t.Start)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %q date", t.Start)
	}
	d2, err := time.Parse("2006-01-02 15:04", t.End)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %q date", t.End)
	}
	d := d2.Sub(d1)
	return d, nil
}

func (db *Db) ListTimesheet(query string) ([]*Time, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("db.listtimesheet error: %s", err)
	}
	defer rows.Close()

	var times []*Time
	for rows.Next() {
		var t Time
		err := rows.Scan(
			&t.Start,
			&t.End,
			&t.Income,
			&t.TotalTime,
		)
		if err != nil {
			return nil, fmt.Errorf("db.listtimesheet error: %s", err)
		}
		t.TotalTimeCalc, err = t.InHours()
		if err != nil {
			return nil, fmt.Errorf("failed to calculate total time: %s", err.Error())
		}
		times = append(times, &t)
	}
	return times, nil
}

type ContractCalendar struct {
	Calendar       cal.BusinessCalendar
	HoursRemaining time.Duration
	DaysRemaining  int
	ContractEnd    time.Time
}

func NewCalendar(end time.Time) *ContractCalendar {
	c := cal.NewBusinessCalendar()
	c.Holidays = au.HolidaysACT
	c.SetWorkHours(8*time.Hour, 17*time.Hour+30*time.Minute) // 8h30m
	return &ContractCalendar{
		ContractEnd: end,
		Calendar:    *c,
	}
}

// DaysLeft calculates total possible working days remaining until contract end
func (c *ContractCalendar) DaysLeft() int {
	return c.Calendar.WorkdaysInRange(time.Now(), c.ContractEnd)
}

// HoursLeft calculates hours remaining in contract
func (c *ContractCalendar) HoursLeft() time.Duration {
	return c.Calendar.WorkHoursInRange(time.Now(), c.ContractEnd)
}

// MeanDaily calculates average daily hours needed to finish contract at zero hours
func (c *ContractCalendar) MeanDaily(contractHours, totalSeconds int) (float64, error) {
	ch, err := ContractHoursToNow(totalSeconds)
	if err != nil {
		return 0, err
	}
	return (float64(contractHours) - ch.Hours()) / float64(c.DaysLeft()), nil
}

// DaysLeftThisMonth calculates number of working days this month
func (c *ContractCalendar) DaysLeftThisMonth() int {
	remainder := cal.DayStart(cal.MonthStart(time.Now().AddDate(0, 1, 0)))
	return c.Calendar.WorkdaysInRange(time.Now(), remainder)
}

// below is todo later

// EndOfMonthTrigger creates an event based on the time being the end of this
// calendar month. It takes in an Event and publishes it using a Publisher.
func (c *ContractCalendar) EndOfMonthTrigger(e Event) (*Publisher, error) {
	if c.Calendar.WorkdaysRemain(time.Now()) == 0 {
		return &Publisher{}, nil
	}
	return nil, nil
}

type Event struct {
	message any
	to      string
	from    string
}
type Publisher struct {
	carrier Carrier
	data    Event
}

const (
	Email Carrier = iota
	Telegram
)

type Carrier int64

func (c Carrier) String() string {
	switch c {
	case Email:
		return "email"
	case Telegram:
		return "telegram"
	default:
		return "email"
	}
}
