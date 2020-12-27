$(function(){

    const inputField = document.getElementById('q');
    const samplePopup = document.getElementById('samples');

    function showError(error) {
        let html = "";

        for (let i = 0; i < error.length; i++) {
            html += error[i] + "<br>";
        }

        document.getElementById('error-box').innerHTML = html;
    }

    function showAnswer(answer) {
        document.getElementById('answer-box').innerHTML = answer;
    }

    function clearInput() {
        document.getElementById('q').value = "";
    }

    function showProductions(productions) {

        let html = '';

        for (let key in productions) {
            let production = productions[key];

            let matches = production.match(/([^:]+)/);
            let name = matches[1];
            let value = production.substr(name.length + 1)
                .replace(/&/g, "&amp;")
                .replace(/</g, "&lt;")
                .replace(/>/g, "&gt;")
                .replace(/"/g, "&quot;")
                .replace(/'/g, "&#039;")
                .replace("\n", "<br>");

            html += "<h2>" + name + "</h2>";
            html += "<p>" + value + "</p>";
        }

        document.getElementById('production-box').innerHTML = html;
    }

    $('#close').click(function() {
        let popup = document.getElementById('popup');
        popup.style.display = "none";
    });

    $('#samples a').click(function(event){
        event.preventDefault();
        $('#q').val(this.innerHTML);
        $('#samples').hide();
    });

    $('#show-samples').click(function (event) {
        event.preventDefault();
        $('#samples').show();
    });

    $('#f').submit(function(){

        postQuestion($('#q').val());
        return false;
    });

    let samples = document.querySelectorAll('#samples li');
    for (let i = 0; i < samples.length; i++) {
        samples[i].onclick = function (element) {
            let li = element.currentTarget;
            inputField.value = li.innerHTML;
            samplePopup.style.display = "none";
        }
    }

    function postQuestion(question) {
        $.ajax({
            url: 'ajax-answer.php',
            data: { format: "json", query: question },
            dataType: 'json',
            type: 'GET',
            success: function (data) {

                if (data.OptionKeys.length === 0) {
                    showAnswer(data.Answer);
                    clearInput();
                    log(question, data.Answer)
                } else {
                    showAnswer("");
                }
                showError(data.ErrorLines);
                showProductions(data.Productions);
                showOptions(data.Answer, data.OptionKeys, data.OptionValues);
            },
            error: function (request, status, error) {
                showError(error)
            }
        });
    }

    function showOptions(answer, optionKeys, optionValues) {
        let html = "<ol>";
        let showOptions = optionKeys.length > 0;

        for (let i = 0; i < optionKeys.length; i++) {
            html += "<li><a href='" + optionKeys[i] + "'>" + optionValues[i] + "</a></li>";
        }

        html += "</ol>"

        let popup = document.getElementById('popup');
        let optionsBox = document.getElementById('options-box');
        let optionsHeader = document.getElementById('options-header');

        optionsHeader.innerHTML = answer;
        optionsBox.innerHTML = html;

        popup.style.display = showOptions ? "block" : "none";

        $('#options-box a').click(function (event) {
            event.preventDefault();
            postQuestion($(this).attr('href'));
        });
    }

    function log(question, answer) {
        let html = "";

        html += "<div><h3>" + question + "</h3></div>";
        html += "<div>" + answer + "</div>";

        let log = document.getElementById("log");
        log.innerHTML = html + log.innerHTML;
    }
});
