package hours

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/sarulabs/di"
)

type HoursTemplate struct {
	db      *sql.DB
	Id      int
	Ticket  string
	Title   string
	Comment string
	Hours   float64
}

type Hours struct {
	db      *sql.DB
	Id      int
	Date    string
	Ticket  string
	Title   string
	Comment string
	Hours   float64
}

type HoursCompat struct {
	Date   string
	Ticket string
	Hours  float64
}

type HoursCalc struct {
	Date  string
	Hours float64
}

func AppendHoursRecordHeader(w table.Writer) {
	w.AppendHeader(table.Row{"Id", "Date", "Ticket", "Title", "Comment", "Hours"})
	w.AppendSeparator()
}

func AppendHoursTemplatesRecordHeader(w table.Writer) {
	w.AppendHeader(table.Row{"Id", "Ticket", "Title", "Comment", "Hours"})
	w.AppendSeparator()
}

func AppendHoursCompatRecordHeader(w table.Writer) {
	w.AppendHeader(table.Row{"Date", "Ticket", "Hours"})
	w.AppendSeparator()
}

func AppendHoursCalcRecordHeader(w table.Writer) {
	w.AppendHeader(table.Row{"Date", "Hours"})
	w.AppendSeparator()
}

func (r Hours) AppendRow(w table.Writer) {
	w.AppendRow([]interface{}{r.Id, r.Date, r.Ticket, r.Title, r.Comment, r.Hours})
	w.AppendSeparator()
}

func (r HoursTemplate) AppendRow(w table.Writer) {
	w.AppendRow([]interface{}{r.Id, r.Ticket, r.Title, r.Comment, r.Hours})
	w.AppendSeparator()
}

func (r HoursCompat) AppendRow(w table.Writer) {
	w.AppendRow([]interface{}{r.Date, r.Ticket, r.Hours})
	w.AppendSeparator()
}

func (r HoursCalc) AppendRow(w table.Writer) {
	w.AppendRow([]interface{}{r.Date, r.Hours})
	w.AppendSeparator()
}

func (r *HoursTemplate) Scan(rows *sql.Rows) {
	rows.Scan(&r.Id, &r.Ticket, &r.Title, &r.Comment, &r.Hours)
}

func (r *Hours) Scan(rows *sql.Rows) {
	rows.Scan(&r.Id, &r.Date, &r.Ticket, &r.Title, &r.Comment, &r.Hours)
}

func (r *HoursCompat) Scan(rows *sql.Rows) {
	rows.Scan(&r.Date, &r.Ticket, &r.Hours)
}

func (r *HoursCalc) Scan(rows *sql.Rows) {
	rows.Scan(&r.Date, &r.Hours)
}

func (r Hours) Validate() error {
	switch {
	case len([]rune(r.Date)) == 0:
	case len([]rune(r.Ticket)) == 0:
	case r.Hours == 0:
		return fmt.Errorf("HoursRecord Validate: record is invalid!! %a", r)
	}

	return nil
}

func readFromStdinSuggestion(name string, suggestion string) (res string, err error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter %s (%s): ", name, suggestion)

	res, err = reader.ReadString('\n')
	return
}

func readFromStdin(name string) (res string, err error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter %s: ", name)

	res, err = reader.ReadString('\n')
	return
}

func NewHours(cnt di.Container, id int, date string, ticket string, title string, comment string, hours float64) Hours {
	db, err := cnt.SafeGet("db")

	if err != nil {
		panic(err)
	} else {
		return Hours{db.(*sql.DB), id, date, ticket, title, comment, hours}
	}
}

func NewHoursTemplate(cnt di.Container, id int, ticket string, title string, comment string, hours float64) HoursTemplate {
	db, err := cnt.SafeGet("db")

	if err != nil {
		panic(err)
	} else {
		return HoursTemplate{db.(*sql.DB), id, ticket, title, comment, hours}
	}
}

func (r Hours) FillOut() Hours {
	if r.Date == "" {
		text, err := readFromStdinSuggestion("Date", "leave empty for today")

		if err != nil {
			panic(err)
		} else {
			sDate := strings.TrimSuffix(text, "\n")

			if sDate == "" {
				t := time.Now()
				r.Date = fmt.Sprintf("%d-%d-%d", t.Year(), t.Month(), t.Day())
			} else {
				r.Date = sDate
			}
		}
	}

	if r.Ticket == "" {
		text, err := readFromStdin("Ticket")

		if err != nil {
			panic(err)
		} else {
			r.Ticket = strings.TrimSuffix(text, "\n")
		}
	}

	if r.Title == "" {
		text, err := readFromStdin("Title")

		if err != nil {
			panic(err)
		} else {
			r.Title = strings.TrimSuffix(text, "\n")
		}
	}

	if r.Comment == "" {
		text, err := readFromStdinSuggestion("Comment", "can be left empty")

		if err != nil {
			panic(err)
		} else {
			r.Comment = strings.TrimSuffix(text, "\n")
		}
	}

	if r.Hours == 0 {
		text, err := readFromStdin("Hours")

		if err != nil {
			panic(err)
		} else {
			fHours, err := strconv.ParseFloat(strings.TrimSuffix(text, "\n"), 64)
			if err != nil {
				panic(err)
			}

			r.Hours = fHours
		}
	}

	return r
}

func (r Hours) Insert() (sql.Result, error) {
	statement, _ := r.db.Prepare(
		"INSERT INTO hours (date, ticket, title, comment, hours) VALUES (?, ?, ?, ?, ?)",
	)
	defer statement.Close()

	return statement.Exec(r.Date, r.Ticket, r.Title, r.Comment, r.Hours)
}

func (r Hours) Update() (sql.Result, error) {
	statement, _ := r.db.Prepare(
		"UPDATE hours SET date = ?, ticket = ?, title = ?, comment = ?, hours = ? WHERE id = ?",
	)
	defer statement.Close()

	return statement.Exec(r.Date, r.Ticket, r.Title, r.Comment, r.Hours, r.Id)
}

func (r HoursTemplate) Insert() (sql.Result, error) {
	statement, _ := r.db.Prepare(
		"INSERT INTO hours_templates (ticket, title, comment, hours) VALUES (?, ?, ?, ?)",
	)
	defer statement.Close()

	return statement.Exec(r.Ticket, r.Title, r.Comment, r.Hours)
}

func (r HoursTemplate) Update() (sql.Result, error) {
	statement, _ := r.db.Prepare(
		"UPDATE hours_templates SET ticket = ?, title = ?, comment = ?, hours = ? WHERE id = ?",
	)
	defer statement.Close()

	return statement.Exec(r.Ticket, r.Title, r.Comment, r.Hours, r.Id)
}
