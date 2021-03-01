"use strict"

$(function () {
    $("#add").click(function () {
        let modalSel = "#modal-container";

        $(modalSel).load("templates/add", function () {
            $("#modal-save-template").click(function () {
                $.post("templates/add", {
                    id: $("#modal-hour-id").val(),
                    ticket: $("#modal-ticket").val(),
                    title: $("#modal-title").val(),
                    comment: $("#modal-comment").val(),
                    hours: $("#modal-hours").val()
                })
                .then(function() {
                    window.location.reload();
                })
                .fail(function() {
                    alert("There has been an error, check the webserver logs")
                });
            });

            $(modalSel).modal("show");
        });
    });
});