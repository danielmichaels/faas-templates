package handlers

var (
	ContractStart = `SELECT date1,date2,amount,ROUND((JULIANDAY(date2) - JULIANDAY(date1))*86400)
    AS diff
    FROM times
    WHERE date1 between '2022-11-30' AND JULIANDAY('now')`

	OldContract = `SELECT date1,date2,amount,ROUND((JULIANDAY(date2) - JULIANDAY(date1))*86400)
    AS diff
    FROM times
    WHERE date1 between '2020-11-30' AND '2022-11-29'`

	LastSevenDays = `SELECT date1,date2,amount,ROUND((JULIANDAY(date2)-JULIANDAY(date1))*86400) 
    AS diff 
	FROM times 
	WHERE date1 > datetime('now', '-7 days');`

	LastCalendarMonth = `SELECT date1,date2,amount,ROUND((JULIANDAY(date2)-JULIANDAY(date1))*86400) 
    AS diff 
	FROM times 
	WHERE date1
	BETWEEN datetime('now', 'start of month') AND datetime('now', 'localtime');`
)
