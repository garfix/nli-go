$(function(){

    function showError(error) {
        document.getElementById('error-box').innerHTML = error;
    }

    function showAnswer(answer) {
        document.getElementById('answer-box').innerHTML = answer;
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

                if (!success) {
                    showError(errorLines)
                } else {

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
                var success = data.Success;
                var errorLines = data.ErrorLines;

                if (!success) {
                    showError(errorLines)
                } else {
                    showAnswer(answers[0])
                }
            },
            error: function (request, status, error) {
                showError(error)
            }
        });

        return false;
    })
});
