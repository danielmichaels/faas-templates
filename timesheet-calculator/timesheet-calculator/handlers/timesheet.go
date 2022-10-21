package handlers

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"time"
)

var (
	ContractStart = `SELECT date1,date2,amount,ROUND((JULIANDAY(date2) - JULIANDAY(date1))*86400)
    AS diff
    FROM times
    WHERE date1 between '2022-11-30' AND JULIANDAY('now')`
	OldContract = `SELECT date1,date2,amount,ROUND((JULIANDAY(date2) - JULIANDAY(date1))*86400)
    AS diff
    FROM times
    WHERE date1 between '2020-11-30' AND '2022-11-29'`
	OldContractEnd = time.Date(2022, 11, 30, 0, 0, 0, 00, time.UTC)
	NewContractEnd = time.Date(2024, 07, 00, 0, 0, 0, 00, time.UTC)
)

func ContractHoursToNow(secs int) (time.Duration, error) {
	h, err := time.ParseDuration(fmt.Sprintf("%ds", secs))
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration: %w", err)
	}
	return h, nil
}

// searchTimesheetBackward searches over a date time range of fourteen (14) days
// to find a matching database. Database list is sorted with the most recent
// first.
func searchTimesheetBackward(ts []string) (string, error) {
	searchableDates := []time.Time{
		time.Now(),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -2),
		time.Now().AddDate(0, 0, -3),
		time.Now().AddDate(0, 0, -4),
		time.Now().AddDate(0, 0, -5),
		time.Now().AddDate(0, 0, -6),
		time.Now().AddDate(0, 0, -7),
		time.Now().AddDate(0, 0, -8),
		time.Now().AddDate(0, 0, -9),
		time.Now().AddDate(0, 0, -10),
		time.Now().AddDate(0, 0, -11),
		time.Now().AddDate(0, 0, -12),
		time.Now().AddDate(0, 0, -13),
	}
	sort.Sort(sort.Reverse(sort.StringSlice(ts)))
	for _, d := range searchableDates {
		tFmt := fmt.Sprintf("%d_%02d_%02d", d.Year(), d.Month(), d.Day())
		ref := fmt.Sprintf("timesheets/timesheet_%s_auto_database.db", tFmt)
		for _, t := range ts {
			if ref == t {
				log.Println("most recent database", ref)
				return ref, nil
			}
		}
	}
	return "", errors.New("no database found in the last 14 days")
}

// MostRecentTimesheet returns the most recent timesheet in the DataStore.
func MostRecentTimesheet(ts []string) (string, error) {
	return searchTimesheetBackward(ts)
}
