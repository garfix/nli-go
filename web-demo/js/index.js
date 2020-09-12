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

                showError(data.ErrorLines);
                showAnswer(data.Answer);
                showProductions(data.Productions);
                showOptions(data.OptionKeys, data.OptionValues);

            },
            error: function (request, status, error) {
                showError(error)
            }
        });
    }

    function showOptions(optionKeys, optionValues) {
        var html = "";

        if (optionKeys) {

            for (var i = 0; i < optionKeys.length; i++) {
                html += "<a href='" + optionKeys[i] + "'>" + optionValues[i] + "</a>";
            }

            if (html) {
                html = "<div>" + html + "</div>";
            }
        }

        document.getElementById('options-box').innerHTML = html;

        $('#options-box a').click(function (event) {
            event.preventDefault();
            postQuestion($(this).attr('href'));
        });

    }
});
