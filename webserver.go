package main

import (
	"encoding/json"
	"log"
	"strconv"

	records "internal/hoursmngr/records"

	"github.com/kataras/iris/v12"
	"github.com/sarulabs/di"
)

const baseTitle = "Hours Manager by λ"

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
		"title":   baseTitle + " - Index",
		"hours":   hours,
		"summary": hoursSummary,
	}

	ctx.ViewLayout("layouts/main")
	ctx.View("list", data)
}

func listHoursTemplates(ctx iris.Context, cnt di.Container) {
	var (
		hoursTemplates []records.HoursTemplate
	)

	rawHoursTemplatesRows, err := showTemplates(cnt)
	if err != nil {
		log.Fatal("listHoursTemplates error: " + err.Error())
	}
	defer rawHoursTemplatesRows.Close()

	for rawHoursTemplatesRows.Next() {
		var row records.HoursTemplate

		row.Scan(rawHoursTemplatesRows)
		hoursTemplates = append(hoursTemplates, row)
	}

	data := iris.Map{
		"title":     baseTitle + " - Templates Editor",
		"templates": hoursTemplates,
	}

	ctx.ViewLayout("layouts/main")
	ctx.View("listTemplates", data)
}

func addHourInterface(ctx iris.Context, cnt di.Container) {
	data := iris.Map{
		"modalTitle": "Add Hour",
	}

	ctx.View("edit", data)
}

func dupHourInterface(ctx iris.Context, cnt di.Container) {
	rid := ctx.PostValue("id")
	id, err := strconv.Atoi(rid)
	if err != nil {
		log.Println("addHour insert id validate error: " + err.Error())
		ctx.SetErr(err)
		return
	}

	hour, err := getHour(cnt, id)

	data := iris.Map{
		"modalTitle": "Edit Hour",
		"hour":       hour,
	}

	ctx.View("duplicate", data)
}

func addTemplateInterface(ctx iris.Context, cnt di.Container) {
	data := iris.Map{
		"modalTitle": "Add Hour Template",
	}

	ctx.View("templateEdit", data)
}

func editHour(ctx iris.Context, cnt di.Container) {
	rid := ctx.PostValue("id")
	date := ctx.PostValue("date")
	ticket := ctx.PostValue("ticket")
	title := ctx.PostValue("title")
	comment := ctx.PostValue("comment")
	rhours := ctx.PostValue("hours")

	hours, err := strconv.ParseFloat(rhours, 32)
	if err != nil {
		log.Println("addHour insert hours validate error: " + err.Error())
		ctx.JSON(iris.Map{
			"error": true,
		})
		return
	}

	record := records.NewHours(
		cnt,
		0,
		date,
		ticket,
		title,
		comment,
		hours,
	)

	verr := record.Validate()

	if verr != nil {
		log.Println("addHour insert validate error: " + verr.Error())
		ctx.JSON(iris.Map{
			"error": true,
		})
		return
	}

	if rid != "" {
		id, err := strconv.Atoi(rid)
		if err != nil {
			log.Println("addHour insert id validate error: " + err.Error())
			ctx.JSON(iris.Map{
				"error": true,
			})
			return
		}

		record.Id = id
		_, err = record.Update()

		if err != nil {
			log.Println("addHour update error: " + err.Error())
			ctx.JSON(iris.Map{
				"error": true,
			})
			return
		}
	} else {
		_, err := record.Insert()

		if err != nil {
			log.Println("editHour insert error: " + err.Error())
			ctx.JSON(iris.Map{
				"error": true,
			})
			return
		}
	}

	ctx.JSON(iris.Map{})
}

func editTemplate(ctx iris.Context, cnt di.Container) {
	rid := ctx.PostValue("id")
	ticket := ctx.PostValue("ticket")
	title := ctx.PostValue("title")
	comment := ctx.PostValue("comment")
	rhours := ctx.PostValue("hours")

	hours, err := strconv.ParseFloat(rhours, 32)
	if err != nil {
		log.Println("addHour insert hours validate error: " + err.Error())
		ctx.JSON(iris.Map{
			"error": true,
		})
		return
	}

	record := records.NewHoursTemplate(
		cnt,
		0,
		ticket,
		title,
		comment,
		hours,
	)

	if rid != "" {
		id, err := strconv.Atoi(rid)
		if err != nil {
			log.Println("addHour insert id validate error: " + err.Error())
			ctx.JSON(iris.Map{
				"error": true,
			})
			return
		}

		record.Id = id
		_, err = record.Update()

		if err != nil {
			log.Println("addHour update error: " + err.Error())
			ctx.JSON(iris.Map{
				"error": true,
			})
			return
		}
	} else {
		_, err := record.Insert()

		if err != nil {
			log.Println("editHour insert error: " + err.Error())
			ctx.JSON(iris.Map{
				"error": true,
			})
			return
		}
	}

	ctx.JSON(iris.Map{})
}

func deleteHours(ctx iris.Context, cnt di.Container) {
	rid := ctx.PostValueDefault("id", "")

	if rid != "" {
		id, err := strconv.Atoi(rid)
		if err != nil {
			log.Println("delete hour id validate error: " + err.Error())
			ctx.JSON(iris.Map{
				"error": true,
			})
			return
		}

		_, err = clearById(cnt, id)
		if err != nil {
			log.Println("delete hour error: " + err.Error())
			ctx.JSON(iris.Map{
				"error": true,
			})
			return
		}
	} else {
		//TODO: debug this branch
		var ids []int

		rids := ctx.PostValueDefault("ids", "[]")
		err := json.Unmarshal([]byte(rids), &ids)

		if err != nil {
			log.Println("delete hour ids unmarshal error: " + err.Error())
			ctx.JSON(iris.Map{
				"error": true,
			})
			return
		}

		for _, id := range ids {
			_, err = clearById(cnt, id)
			if err != nil {
				log.Println("delete hour error: " + err.Error())
				ctx.JSON(iris.Map{
					"error": true,
				})
				return
			}
		}

	}

	ctx.JSON(iris.Map{})
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
		hoursAPI.Get("/add", func(ctx iris.Context) { addHourInterface(ctx, cnt) })

		hoursAPI.Post("/edit", func(ctx iris.Context) { editHour(ctx, cnt) })
		hoursAPI.Post("/duplicate", func(ctx iris.Context) { dupHourInterface(ctx, cnt) })
		hoursAPI.Post("/delete", func(ctx iris.Context) { deleteHours(ctx, cnt) })

		/*
			hoursAPI.Get("/addTemplate", addHourWithTemplateInterface)
		*/
	}

	hoursTemplatesAPI := app.Party("/templates")
	{
		hoursTemplatesAPI.Get("/", func(ctx iris.Context) { listHoursTemplates(ctx, cnt) })
		hoursTemplatesAPI.Get("/add", func(ctx iris.Context) { addTemplateInterface(ctx, cnt) })

		hoursTemplatesAPI.Post("/edit", func(ctx iris.Context) { editTemplate(ctx, cnt) })
		/*
			hoursTemplatesAPI.Post("/delete", deleteTemplate)
		*/
	}

	app.Get("/", func(ctx iris.Context) {
		r := ctx.Request()
		r.URL.Path = "/hours"

		ctx.Exec("GET", "/hours")
	})
	app.Listen(":8000")
}
