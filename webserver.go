package main

import (
	"log"

	records "internal/hoursmngr/records"

	"github.com/kataras/iris/v12"
	"github.com/sarulabs/di"
)

func listHours(ctx iris.Context, cnt di.Container) {
	var (
		hours        []records.Hours
		hoursSummary []records.HoursCompat
	)

	rawHoursRows, err := show(cnt)
	if err != nil {
		log.Fatal("listHours error: " + err.Error())
	}
	defer rawHoursRows.Close()

	rawHoursSummaryRows, err := compat(cnt)
	if err != nil {
		log.Fatal("listHours error: " + err.Error())
	}
	defer rawHoursSummaryRows.Close()

	for rawHoursRows.Next() {
		var row records.Hours

		row.Scan(rawHoursRows)
		hours = append(hours, row)
	}

	for rawHoursSummaryRows.Next() {
		var row records.HoursCompat

		row.Scan(rawHoursSummaryRows)
		hoursSummary = append(hoursSummary, row)
	}

	data := iris.Map{
		"hours":   hours,
		"summary": hoursSummary,
	}

	ctx.ViewLayout("layouts/main")
	ctx.View("list", data)
}

func serve(cnt di.Container) {
	app := iris.New()
	app.HandleDir("/pub", iris.Dir("./pub"))

	tmpl := iris.HTML("./views", ".html")
	tmpl.Delims("{{", "}}")
	tmpl.Reload(true)
	app.RegisterView(tmpl)

	hoursAPI := app.Party("/hours")
	{
		hoursAPI.Get("/", func(ctx iris.Context) { listHours(ctx, cnt) })
		/*
			hoursAPI.Get("/add", addHourInterface)
			hoursAPI.Get("/addTemplate", addHourWithTemplateInterface)

			hoursAPI.Post("/add", addHour)
			hoursAPI.Post("/delete", deleteHours)
		*/
	}
	/*
		hoursTemplatesAPI := app.Party("/hoursTpl")
		{

				hoursTemplatesAPI.Get("/", listTemplates)
				hoursTemplatesAPI.Get("/add", addTemplateInterface)

				hoursTemplatesAPI.Post("/add", addTemplate)
				hoursTemplatesAPI.Post("/delete", deleteTemplate)

		}
	*/

	app.Get("/", func(ctx iris.Context) {
		r := ctx.Request()
		r.URL.Path = "/hours"

		ctx.Exec("GET", "/hours")
	})
	app.Listen(":8000")
}
