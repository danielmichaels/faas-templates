package function

import (
	"fmt"
	"handler/function/handlers"
	"log"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {

	b2, err := handlers.NewB2Client()
	if err != nil {
		log.Fatalln("b2 failed to initialise error:", err)
	}
	list, err := b2.List()
	if err != nil {
		log.Fatalln("b2 failed to list objects error:", err)
	}
	fmt.Println(list)
	err = b2.Download("timesheets/timesheet_2022_10_13_auto_database.db", "/tmp/timesheet.db")
	if err != nil {
		log.Fatalln("b2 failed to list objects error:", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Body: %s", map[string]any{"list": list})))
}
