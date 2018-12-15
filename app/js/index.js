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
