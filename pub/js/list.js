"use strict"

$(function () {
    let modalSel = "#modal-container";

    function saveHour() {
        let payload = {
            id: $("#modal-hour-id").val(),
            date: $("#modal-date").val(),
            ticket: $("#modal-ticket").val(),
            title: $("#modal-title").val(),
            comment: $("#modal-comment").val(),
            hours: $("#modal-hours").val()
        };

        if (!payload.date || !payload.date === "") {
            alert("Please specify a date!");
            $("#modal-date").focus()
            return;
        }

        if (!payload.ticket || !payload.ticket === "") {
            alert("Please specify a ticket!");
            $("#modal-ticket").focus()
            return;
        }

        if (!payload.hours || !payload.hours === "") {
            alert("Please specify hours!");
            $("#modal-hours").focus()
            return;
        }

        $.post("hours/edit", payload)
            .then(function () {
                window.location.reload();
            })
            .fail(function () {
                alert("There has been an error, check the webserver logs")
            });

    }

    $("#add").click(function () {
        $(modalSel).load("hours/add", function () {
            $("#modal-save-hour").click(saveHour);

            $(modalSel).modal("show");
        });
    });

    $("#addTemplate").click(function () {
        $(modalSel).load("hours/addTemplate", function () {
            $("#modal-use-template").click(function () {
                let templateId = $("#selector").val();

                if (templateId != -1) {
                    $(modalSel).load("hours/addFromTemplate", { id: templateId }, function () {
                        $("#modal-save-hour").click(saveHour);

                        $(modalSel).modal("show");
                    });
                } else {
                    alert("Please choose a valid template");
                    $("#selector").focus();
                }
            });

            $(modalSel).modal("show");
        });
    });

    $("#delete").click(function () {
        let ids = $('input[type="checkbox"][data-id]:checked').get().map((el) => $(el).data("id"));

        if ((ids != []) && (confirm("Are you sure you want to delete the selected rows?"))) {
            $.post("hours/delete", { ids: ids })
                .then(function () {
                    window.location.reload();
                })
                .fail(function () {
                    alert("There has been an error, check the webserver logs")
                });
        }
    });

    $(".dup[data-id]").click(function () {
        let id = $(this).data("id");

        $(modalSel).load("hours/duplicate", { id: id }, function () {
            $("#modal-save-hour").click(saveHour);

            $(modalSel).modal("show");
        });
    });

    $(".delete[data-id]").click(function () {
        let id = $(this).data("id");

        if (confirm("Are you sure you want to delete row #" + id + " ?")) {
            $.post("hours/delete", { id: id })
                .then(function () {
                    window.location.reload();
                })
                .fail(function () {
                    alert("There has been an error, check the webserver logs")
                });
        }
    });
});