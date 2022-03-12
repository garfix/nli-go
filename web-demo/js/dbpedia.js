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

        let question = $('#q').val()

        sendRequest([{
            predicate: 'go_tell',
            arguments: [
                {
                    type: 'string',
                    value: question
                }
            ]
        }]);
        return false;
    });

    function sendRequest(request) {
        $.ajax({
            url: 'ajax-answer.php',
            data: { format: "json", request: JSON.stringify(request), app: "dbpedia" },
            dataType: 'json',
            type: 'GET',
            success: function (data) {

                clearOptionsPopup()
                showProductions(data.Productions);

                if (data.Success) {
                    processResponse(data.Message);
                    showError('');
                } else {
                    showError(data.ErrorLines);
                    showAnswer('')
                }
            },
            error: function (request, status, error) {
                showError(error)
            }
        });
    }

    function processResponse(response) {
        let asserts = [];
        let assert;
        let maxDuration = 0;

        for (let i = 0; i < response.length; i++) {
            let relation = response[i];
            switch (relation.predicate) {
                case 'go_print':
                    assert = print(relation)
                    asserts.push(assert)
                    break;
                case 'go_user_select':
                    showOptionsPopup(relation)
                    break;
            }
        }

        if (asserts.length > 0) {
            window.setTimeout(function (){
                sendRequest(asserts)
            }, maxDuration);
        }
    }

    function print(relation) {
        let answer = relation.arguments[1].value;
        showAnswer(answer)
        return getAssert(relation);
    }

    function getAssert(assertion) {
        return {
            predicate: 'go_assert',
            arguments: [
                {
                    "type": "relation-set",
                    "set": [assertion]
                }
            ]
        }
    }

    function clearOptionsPopup() {
        document.getElementById('options-box').innerHTML = '';
    }

    function showOptionsPopup(relation) {

        let options = relation.arguments[1].list;

        let html = "<ol>";
        for (let i = 0; i < options.length; i++) {
            let argument = options[i];
            html += "<li><a href='" + i + "'>" + argument.value + "</a></li>";
        }
        html += "</ol>"

        document.getElementById('options-box').innerHTML = html;

        $('#options-box a').click(function (event) {
            event.preventDefault();

            relation.arguments[2].type = "string"
            relation.arguments[2].value = event.currentTarget.getAttribute('href');
            let assert = getAssert(relation)
            sendRequest([assert])
        });
    }
});
