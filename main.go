package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sarulabs/di"

	records "internal/hoursmngr/records"
)

func doInsert(record records.Hours) {
	record = record.FillOut()
	verr := record.Validate()

	if verr != nil {
		log.Fatal("doInsert validate error: " + verr.Error())
	}

	_, err := record.Insert()

	if err != nil {
		log.Fatal("doInsert error: " + err.Error())
	}
}

func insertDailyPauseRecord(cnt di.Container) (sql.Result, error) {
	t := time.Now()
	fmtDate := fmt.Sprintf("%d-%d-%d", t.Year(), t.Month(), t.Day())

	record := records.NewHours(cnt, 0, fmtDate, "<REDACTED>", "Pause", "", 0.5)

	return record.Insert()
}

func doPause(cnt di.Container) {
	_, err := insertDailyPauseRecord(cnt)

	if err != nil {
		log.Fatal("doPause error: " + err.Error())
	}
}

func clearById(cnt di.Container, id int) (sql.Result, error) {
	db, err := cnt.SafeGet("db")

	if err != nil {
		log.Fatal("show di error: " + err.Error())
		return nil, err
	} else {
		statement, _ := db.(*sql.DB).Prepare("DELETE FROM hours WHERE id = ?")
		defer statement.Close()

		return statement.Exec(id)
	}
}

func clearDate(cnt di.Container, date string) (sql.Result, error) {
	db, err := cnt.SafeGet("db")

	if err != nil {
		log.Fatal("show di error: " + err.Error())
		return nil, err
	} else {
		statement, _ := db.(*sql.DB).Prepare("DELETE FROM hours WHERE date = ?")
		defer statement.Close()

		return statement.Exec(date)
	}
}

func doClear(cnt di.Container, date string) {
	if date == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter date: ")
		text, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal("doClear error: " + err.Error())
		} else {
			date = strings.TrimSuffix(text, "\n")
		}
	}

	_, err := clearDate(cnt, date)

	if err != nil {
		log.Fatal("doClear error: " + err.Error())
	}
}

func calc(cnt di.Container) (*sql.Rows, error) {
	db, err := cnt.SafeGet("db")

	if err != nil {
		log.Fatal("calc di error: " + err.Error())
		return nil, err
	} else {
		return db.(*sql.DB).Query(
			"SELECT date, SUM(hours) AS hours FROM hours GROUP BY date ORDER BY date",
		)
	}
}

func doCalc(cnt di.Container) {
	rows, err := calc(cnt)
	defer rows.Close()

	if err != nil {
		log.Fatal("doCalc error: " + err.Error())
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	records.AppendHoursCalcRecordHeader(t)

	for rows.Next() {
		var row records.HoursCalc

		row.Scan(rows)
		row.AppendRow(t)
	}

	t.Render()
}

func compat(cnt di.Container) (*sql.Rows, error) {
	db, err := cnt.SafeGet("db")

	if err != nil {
		log.Fatal("compat di error: " + err.Error())
		return nil, err
	} else {
		return db.(*sql.DB).Query("SELECT date, ticket, hours FROM hours ORDER BY date")
	}
}

func doCompat(cnt di.Container) {
	rows, err := compat(cnt)
	defer rows.Close()

	if err != nil {
		log.Fatal("doCompat error: " + err.Error())
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	records.AppendHoursCompatRecordHeader(t)

	for rows.Next() {
		var row records.HoursCompat

		row.Scan(rows)
		row.AppendRow(t)
	}

	t.Render()
}

func show(cnt di.Container) (*sql.Rows, error) {
	db, err := cnt.SafeGet("db")

	if err != nil {
		log.Fatal("show di error: " + err.Error())
		return nil, err
	}

	return db.(*sql.DB).Query("SELECT * FROM hours ORDER BY date")
}

func getHour(cnt di.Container, id int) (records.Hours, error) {
	var hour records.Hours
	db, err := cnt.SafeGet("db")

	if err != nil {
		log.Fatal("show di error: " + err.Error())
		return hour, err
	}

	statement, _ := db.(*sql.DB).Prepare("SELECT * FROM hours WHERE id = ?")
	defer statement.Close()

	row := statement.QueryRow(id)
	err = row.Scan(&hour.Id, &hour.Date, &hour.Ticket, &hour.Title, &hour.Comment, &hour.Hours)

	if err != nil {
		log.Fatal("show di error: " + err.Error())
		return hour, err
	}

	return hour, nil
}

func showTemplates(cnt di.Container) (*sql.Rows, error) {
	db, err := cnt.SafeGet("db")

	if err != nil {
		log.Fatal("show di error: " + err.Error())
		return nil, err
	}

	return db.(*sql.DB).Query("SELECT * FROM hours_templates ORDER BY id")
}

func doShow(cnt di.Container) {
	rows, err := show(cnt)
	defer rows.Close()

	if err != nil {
		log.Fatal("doShow error: " + err.Error())
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	records.AppendHoursRecordHeader(t)

	for rows.Next() {
		var row records.Hours

		row.Scan(rows)
		row.AppendRow(t)
	}

	t.Render()
}

func main() {
	dateStrSuggestion := "(use `date -d strdate +%F` for easier input)"

	var insertCmdDate, clearCmdDate, insertCmdTicket, insertCmdTitle, insertCmdComment string
	var insertCmdHours float64

	insertCmd := flag.NewFlagSet("insert", flag.ExitOnError)
	insertCmd.StringVar(&insertCmdDate, "date", "", "The date of the record "+dateStrSuggestion)
	insertCmd.StringVar(&insertCmdTicket, "ticket", "", "the redmine ticket id (no hashes)")
	insertCmd.StringVar(&insertCmdTitle, "title", "", "the record title")
	insertCmd.StringVar(&insertCmdComment, "comment", "", "the record comment")
	insertCmd.Float64Var(&insertCmdHours, "hours", 0, "the hours amount to report")

	clearCmd := flag.NewFlagSet("clear", flag.ExitOnError)
	clearCmd.StringVar(&clearCmdDate, "date", "", "The date of the record "+dateStrSuggestion)

	fmt.Println()

	builder, _ := createDIbuilder()
	cnt := builder.Build()

	if len(os.Args) < 2 {
		doShow(cnt)
		doCalc(cnt)
		os.Exit(0)
	}

	switch os.Args[1] {
	case "show":
		doShow(cnt)

	case "compat":
		doCompat(cnt)

	case "calc":
		doCalc(cnt)

	case "pause":
		doPause(cnt)
		doShow(cnt)
		doCalc(cnt)

	case "insert":
		insertCmd.Parse(os.Args[2:])

		record := records.NewHours(
			cnt,
			0,
			insertCmdDate,
			insertCmdTicket,
			insertCmdTitle,
			insertCmdComment,
			insertCmdHours,
		)

		doInsert(record)
		doShow(cnt)
		doCalc(cnt)

	case "clear":
		clearCmd.Parse(os.Args[2:])

		doClear(cnt, clearCmdDate)
		doShow(cnt)
		doCalc(cnt)

	case "serve":
		serve(cnt)

	default:
		fmt.Printf("unknown command '%s'\n", os.Args[1])
		os.Exit(1)
	}

	os.Exit(0)
}
