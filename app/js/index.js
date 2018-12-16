$(function(){

    function showError(error) {
        document.getElementById('error-box').innerHTML = error;
    }

    function showAnswer(answer) {
        document.getElementById('answer-box').innerHTML = answer;
    }

    function showProductions(productions) {

        var html = '<table class="productions">';

        for (var key in productions) {
            var production = productions[key];

            var matches = production.match(/([^:]+)/);
            var name = matches[1];
            var value = production.substr(name.length + 1)
                .replace(/&/g, "&amp;")
                .replace(/</g, "&lt;")
                .replace(/>/g, "&gt;")
                .replace(/"/g, "&quot;")
                .replace(/'/g, "&#039;")
                .replace("\n", "<br>");

            var className = name.toLowerCase().replace(' ', '-');

            html += "<tr><td class='production_name'>" + name + "</td>";
            html += "<td class='production_value " + className + "'>" + value + "</td></tr>";
        }

        html += "</table>";

        document.getElementById('production-box').innerHTML = html;
    }

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

    function postQuestion(question) {
        $.ajax({
            url: 'ajax-answer.php',
            data: { format: "json", query: question },
            dataType: 'json',
            type: 'GET',
            success: function (data) {

                var errorLines = data.ErrorLines;
                var productions = data.Productions;
                var answerComponents = processAnswer(data.Value[0]);
                var answer = answerComponents[0];
                var options = answerComponents[1];

                showError(errorLines);
                showAnswer(answer);
                showProductions(productions);
                showOptions(options);

            },
            error: function (request, status, error) {
                showError(error)
            }
        });
    }

    function processAnswer(answer) {

        var text = "";
        var options = [];
        var state = "text";
        var key = "";
        var value = "";

        for (var i = 0; i < answer.length; i++) {
            var c = answer.substr(i, 1);

            if (c === "[") {
                state = "key";

                if (key !== "") {
                    options.push([key, value]);
                    key = "";
                    value = "";
                }

            } else if (c === "]") {
                state = "value";
            } else if (state === "text") {
                text += c;
            } else if (state === "key") {
                key += c;
            } else {
                value += c;
            }
        }

        if (key !== "") {
            options.push([key, value]);
        }

        return [text, options];
    }

    function showOptions(options) {
        var html = "";

        for (var i = 0; i < options.length; i++) {
            var option = options[i];
            html += "<a href='" + option[0] + "'>" + option[1] + "</a>";
        }

        if (html) {
            html = "<div>" + html + "</div>";
        }

        document.getElementById('options-box').innerHTML = html;

        $('#options-box a').click(function (event) {
            event.preventDefault();
            postQuestion($(this).attr('href'));
        });

    }
});
