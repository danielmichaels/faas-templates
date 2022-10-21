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
	timesheetDb = "/tmp/database.db"
	toEmail     = "dan@danielms.site"
)

type Application struct {
	Db     *handlers.Db
	Cal    *handlers.ContractCalendar
	Time   *handlers.TimeModel
	Mailer *handlers.Mailer
	B2     *handlers.Store
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
		Cal:  handlers.NewCalendar(handlers.OldContractEnd),
		Time: handlers.NewTimeModel(db),
		Db:   handlers.NewDB(db),
		Mailer: handlers.NewMailer(
			os.Getenv("email_host"),
			emailPort,
			string(mailUsername),
			string(mailPassword),
			string(mailFrom),
		),
		B2: b2,
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

	avgHourDaily, err := app.Cal.MeanDaily(1920*2, total)
	if err != nil {
		handlers.ServerError(w, r, err)
		return
	}
	data := map[string]any{
		"MonthDaysLeft":     app.Cal.DaysLeftThisMonth(),
		"ContractHoursLeft": app.Cal.HoursLeft(),
		"AverageHours":      avgHourDaily,
		"ContractEnd":       handlers.OldContractEnd.Format("02-01-2006"),
	}
	log.Println("attempting to send email")
	err = app.Mailer.SendDaily(toEmail, data)
	if err != nil {
		log.Println("email failed to send")
		handlers.ServerError(w, r, err)
		return
	}

	log.Println("successfully sent email")
	_ = handlers.WriteJSON(w, http.StatusOK, handlers.Envelope{"status": "OK"}, nil)
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
