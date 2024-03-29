package function

import (
	"context"
	"database/sql"
	"handler/function/handlers"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"strconv"
	"time"

	"net/http"
)

var (
	timesheetDb          = "/tmp/database.db"
	toEmail              = "dan@danielms.site"
	totalContractedHours = 392 * 8 // days * expected hours per day
)

type result struct {
	Daily   bool `json:"daily"`
	Weekly  bool `json:"weekly"`
	Monthly bool `json:"monthly"`
}
type Application struct {
	Db       *handlers.Db
	Cal      *handlers.ContractCalendar
	Time     *handlers.TimeModel
	Mailer   *handlers.Mailer
	B2       *handlers.Store
	emailQty *result
}

func Handle(w http.ResponseWriter, r *http.Request) {
	b2, err := handlers.NewDataStore()
	if err != nil {
		log.Fatalln("b2 failed to initialise error:", err)
	}

	db, err := openDB(timesheetDb)
	if err != nil {
		handlers.ServerError(w, r, err)
		return
	}

	mailPassword, err := handlers.GetSecret("mail-password")
	if err != nil {
		handlers.ServerError(w, r, err)
		return
	}
	mailUsername, err := handlers.GetSecret("mail-username")
	if err != nil {
		handlers.ServerError(w, r, err)
		return
	}
	mailFrom, err := handlers.GetSecret("mail-from")
	if err != nil {
		handlers.ServerError(w, r, err)
		return
	}
	emailPort, err := strconv.Atoi(os.Getenv("email_port"))
	if err != nil {
		handlers.ServerError(w, r, err)
		return
	}

	app := Application{
		Cal:  handlers.NewCalendar(handlers.ContractEnd),
		Time: handlers.NewTimeModel(db),
		Db:   handlers.NewDB(db),
		Mailer: handlers.NewMailer(
			os.Getenv("email_host"),
			emailPort,
			string(mailUsername),
			string(mailPassword),
			string(mailFrom),
		),
		B2:       b2,
		emailQty: &result{},
	}
	list, err := app.B2.Data.List()
	if err != nil {
		handlers.ServerError(w, r, err)
		return
	}

	ts, err := handlers.MostRecentTimesheet(list)
	if err != nil {
		handlers.ErrorResponse(w, r, http.StatusNotFound, err, "database not found")
		return
	}
	err = app.B2.Data.Download(ts, timesheetDb)
	if err != nil {
		handlers.ErrorResponse(w, r, http.StatusNotFound, err, "file not found")
		return
	}
	tlist, err := app.Db.ListTimesheet(handlers.OldContract)
	if err != nil {
		handlers.ErrorResponse(w, r, http.StatusNotFound, err, "file not found")
		return
	}
	var total int
	for _, v := range tlist {
		total += v.TotalTime
	}

	avgHourDaily, err := app.Cal.MeanDaily(totalContractedHours, total)
	if err != nil {
		handlers.ServerError(w, r, err)
		return
	}
	data := map[string]any{
		"MonthDaysLeft":     app.Cal.DaysLeftThisMonth(),
		"ContractHoursLeft": app.Cal.HoursLeft(),
		"AverageHours":      avgHourDaily,
		"ContractEnd":       handlers.ContractEnd.Format("02-01-2006"),
	}
	if app.Cal.Calendar.IsWorkday(time.Now()) {
		log.Println("attempting to send daily email")
		err = app.Mailer.SendDaily(toEmail, data)
		if err != nil {
			log.Println("email failed to send")
			handlers.ServerError(w, r, err)
			return
		}
		log.Println("successfully sent email")
		app.emailQty.Daily = true
	}

	if app.Cal.IsSaturday() {
		log.Println("attempting to send weekly email")
		t, err := app.Db.GetLastWeek(handlers.LastSevenDays)
		if err != nil {
			log.Println("error: failed to determine last weeks information")
			handlers.ErrorResponse(w, r, http.StatusNotFound, err, "not found")
			return
		}

		meanDaily, err := handlers.MeanDuration(t)
		if err != nil {
			log.Println("error: failed to determine last weeks information")
			handlers.ServerError(w, r, err)
			return
		}
		cumulativeHours, err := handlers.CumulativeDuration(t)
		if err != nil {
			log.Println("error: failed to determine cumulative hours")
			handlers.ServerError(w, r, err)
			return
		}
		weeklyData := map[string]any{
			"MonthDaysLeft":     app.Cal.DaysLeftThisMonth(),
			"ContractHoursLeft": app.Cal.HoursLeft(),
			"AverageHours":      avgHourDaily,
			"TotalHours":        cumulativeHours,
			"ContractEnd":       handlers.ContractEnd.Format("02-01-2006"),
			"Income":            handlers.EstimatedIncome(t),
			"MeanDaily":         meanDaily,
			"NumDays":           len(t),
			"Times":             t,
		}

		err = app.Mailer.SendWeekly(toEmail, weeklyData)
		if err != nil {
			log.Println("email failed to send")
			handlers.ServerError(w, r, err)
			return
		}
		log.Println("successfully sent email")
		app.emailQty.Weekly = true
	}

	// End of the month
	if app.Cal.Calendar.WorkdaysRemain(time.Now()) == 0 {
		log.Println("attempting to send month email")
		t, err := app.Db.GetLastMonth(handlers.LastCalendarMonth)
		if err != nil {
			log.Println("error: failed to determine last months information")
			handlers.ErrorResponse(w, r, http.StatusNotFound, err, "not found")
			return
		}

		meanDaily, err := handlers.MeanDuration(t)
		if err != nil {
			log.Println("error: failed to determine last months information")
			handlers.ServerError(w, r, err)
			return
		}
		cumulativeHours, err := handlers.CumulativeDuration(t)
		if err != nil {
			log.Println("error: failed to determine cumulative hours")
			handlers.ServerError(w, r, err)
			return
		}

		monthlyData := map[string]any{
			"ContractHoursLeft": app.Cal.HoursLeft(),
			"AverageHours":      avgHourDaily,
			"ContractEnd":       handlers.ContractEnd.Format("02-01-2006"),
			"Income":            handlers.EstimatedIncome(t),
			"TotalHours":        cumulativeHours,
			"MeanDaily":         meanDaily,
			"Month":             time.Now().Month(),
			"NumDays":           len(t),
			"Times":             t,
		}

		err = app.Mailer.SendMonthly(toEmail, monthlyData)
		if err != nil {
			log.Println("email failed to send")
			handlers.ServerError(w, r, err)
			return
		}
		log.Println("successfully sent email")
		app.emailQty.Monthly = true
	}

	switch {
	case app.emailQty.Daily == true:
		break
	case app.emailQty.Weekly == true:
		break
	case app.emailQty.Monthly == true:
		break
	default:
		log.Println("no emails sent today")
	}

	_ = handlers.WriteJSON(w, http.StatusOK, handlers.Envelope{"status": "OK", "code": 200, "result": app.emailQty}, nil)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
