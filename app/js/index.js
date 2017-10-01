$(function(){

    function showError(error) {
        document.getElementById('error-box').innerHTML = error;
    }

    function showAnswer(answer) {
        document.getElementById('answer-box').innerHTML = answer;
    }

    function showProductions(productions) {

        html = '<table class="productions">';

        for (var key in productions) {
            production = productions[key];

            matches = production.match(/([^:]+)/);
            name = matches[1];
            value = production.substr(name.length + 1)
                .replace(/&/g, "&amp;")
                .replace(/</g, "&lt;")
                .replace(/>/g, "&gt;")
                .replace(/"/g, "&quot;")
                .replace(/'/g, "&#039;")
                .replace("\n", "<br>");

            html += "<tr><td class='production_name'>" + name + "</td>";
            html += "<td class='production_value'>" + value + "</td></tr>";
        }

        html += "</table>";

        document.getElementById('production-box').innerHTML = html;
    }

    function autoCompleteSource(request, response) {

        $.ajax({
            url: 'ajax-suggest.php',
            data: { format: "json", query: $('#q').val() },
            dataType: 'json',
            type: 'GET',
            success: function (data) {

                var suggests = data.Value;
                var success = data.Success;
                var errorLines = data.ErrorLines;

                showError(errorLines);

                if (success) {

                    response($.map(suggests, function (item) {
                        return {
                            label: item,
                            value: item
                        }
                    }));
                }
            },
            error: function (request, status, error) {
                showError(error)
            }
        })
    }

    $('#q').tagit({
        autocomplete: {delay: 0, minLength: 0, source: autoCompleteSource},
        showAutocompleteOnFocus: true,
        removeConfirmation: false,
        caseSensitive: true,
        allowDuplicates: true,
        allowSpaces: false,
        readOnly: false,
        tagLimit: null,
        tabIndex: null,
        placeholderText: "Next word ..."
    });

    $('#f').submit(function(){
        $.ajax({
            url: 'ajax-answer.php',
            data: { format: "json", query: $('#q').val() },
            dataType: 'json',
            type: 'GET',
            success: function (data) {

                var answers = data.Value;
                var errorLines = data.ErrorLines;
                var productions = data.Productions;

                showError(errorLines);
                showAnswer(answers[0]);
                showProductions(productions)

            },
            error: function (request, status, error) {
                showError(error)
            }
        });

        return false;
    })
});
