$(function(){

    const inputField = document.getElementById('q');
    const samplePopup = document.getElementById('samples');
    const productionBoxLeft = document.getElementById('production-box-left');
    const productionBoxRight = document.getElementById('production-box-right');
    const answerBox = document.getElementById('answer-box');
    const errorBox = document.getElementById('error-box');
    const popup = document.getElementById('popup');
    const popupCloseButton = document.getElementById('close');
    const sampleButton = document.getElementById('show-samples');
    const form = document.getElementById('f');
    const optionsBox = document.getElementById('options-box');
    const optionsHeader = document.getElementById('options-header');
    const logBox = document.getElementById("log");
    const monitor = document.getElementById("monitor");
    const resetButton = document.getElementById("reset");

    const displayWidth = Math.min(window.innerWidth, 600);
    const displayHeight = displayWidth / 2;

    let scene = createScene();

    function setup()
    {
        monitor.style.height = displayHeight + "px";

        popupCloseButton.onclick = function() {
            hidePopup()
        };

        sampleButton.onclick = function (event) {
            showPopup()
        };

        resetButton.onclick = function () {
            reset();
        }

        form.onsubmit = function(){
            tell(inputField.value)
            return false;
        };

        document.addEventListener('keydown', function(event){
            if (event.key === 'Escape') {
                hidePopup()
            }
        });

        let samples = document.querySelectorAll('#samples li');
        for (let i = 0; i < samples.length; i++) {
            samples[i].onclick = function (event) {
                let li = event.target;
                inputField.value = li.innerHTML;
                samplePopup.style.display = "none";
            }
        }

        logBox.onclick = function (event) {
            let e = event.target;
            if (e.tagName.toLowerCase() === "h3") {
                inputField.value = e.innerHTML;
            }
        }

        updateScene(true)
    }

    function showPopup() {
        samplePopup.style.display = "block";
    }

    function hidePopup() {
        samplePopup.style.display = "none";
    }

    function showError(error) {
        let html = "";

        if (Array.isArray(error)) {
            for (let i = 0; i < error.length; i++) {
                html += error[i] + "<br>";
            }
        } else {
            html = error;
        }

        errorBox.innerHTML = html;
    }

    function showAnswer(answer) {
        answerBox.innerHTML = answer;
    }

    function clearInput() {
        inputField.value = "";
    }

    function showProductions(productions) {

        let html = {
            'production-box-left': "",
            'production-box-right': ""
        };

        let container = 'production-box-left';

        for (let key in productions) {
            let production = productions[key];

            let matches = production.match(/([^:]+)/);
            let name = matches[1];
            let value = production.substr(name.length + 2);

            if (name === "Named entities") { continue }

            html[container] += "<div class='card'><h2>" + name + "</h2>" + "<pre>" + value + "</pre></div>";

            if (name === "Parse tree") {
                container = 'production-box-right';
            }
        }

        productionBoxLeft.innerHTML = html['production-box-left'];
        productionBoxRight.innerHTML = html['production-box-right'];
    }

    function reset() {
        $.ajax({
            url: 'scene.php',
            data: { format: "json", action: 'reset' },
            dataType: 'json',
            type: 'GET',
            success: function () {
                window.location.reload();
            },
            error: function (request, status, error) {
                showError(error)
            }
        });
    }

    function sendRequest(request) {
        $.ajax({
            url: 'ajax-answer.php',
            data: { format: "json", request: JSON.stringify(request), app: "blocks" },
            dataType: 'json',
            type: 'GET',
            success: function (data) {

                processResponse(data.AnswerStruct)

                // if (data.OptionKeys.length === 0) {
                //     showAnswer(data.Answer);
                //     clearInput();
                //     log(question, data.Answer)
                //     updateScene(false)
                // } else {
                //     showAnswer("");
                // }
                showError(data.ErrorLines);
                showProductions(data.Productions);
                //showOptions(data.Answer, data.OptionKeys, data.OptionValues);
            },
            error: function (request, status, error) {
                showError(error)
            }
        });
    }

    function processResponse(response)
    {
        for (let i = 0; i < response.length; i++) {
            let relation = response[i];
            switch (relation.predicate) {
                case 'dom_move_object':
                    moveObject(objectId, response[1], response[2], response[3])
                    break;
                case 'go_print':
                    print(relation)
                    break;
                case 'go_user_select':
                    initUserSelect(response[1])
                    break;
            }
        }
    }

    function print(relation) {
        showAnswer(relation.arguments[1].value)
        assert(relation);
    }

    function tell(input) {
        sendRequest({
            positive: true,
            predicate: 'go_tell',
            arguments: [
                {
                    type: 'string',
                    value: input
                }
            ]
        });
    }

    function assert(assertion) {
        sendRequest({
            positive: true,
            predicate: 'go_assert',
            arguments: [
                {
                    "type": "relation-set",
                    "set": [assertion]
                }
            ]
        });
    }

    function showOptions(answer, optionKeys, optionValues) {
        let html = "<ol>";
        let showOptions = optionKeys.length > 0;

        for (let i = 0; i < optionKeys.length; i++) {
            html += "<li><a href='" + optionKeys[i] + "'>" + optionValues[i] + "</a></li>";
        }

        html += "</ol>"

        optionsHeader.innerHTML = answer;
        optionsBox.innerHTML = html;

        popup.style.display = showOptions ? "block" : "none";

        let aTags = optionsBox.querySelectorAll('a');
        for (let i = 0; i < aTags.length; i++) {
            aTags[i].onclick = function (event) {
                event.preventDefault();
                sendRequest(event.currentTarget.getAttribute('href'));
            };
        }
    }

    function log(question, answer) {
        let html = "";

        html += "<div class='prev-question'><h3>" + question + "</h3></div>";
        html += "<div class='prev-answer'>" + answer + "</div>";

        logBox.innerHTML = html + logBox.innerHTML;
    }

    function updateScene(initial)
    {
        $.ajax({
            url: 'scene.php',
            data: { format: "json", action: "state" },
            dataType: 'json',
            type: 'GET',
            success: function (data) {
                if (initial) {
                    scene.build(monitor, data, displayWidth, displayHeight)
                } else {
                    scene.update(data)
                }
            },
            error: function (request, status, error) {
                showError(error)
            }
        });
    }

    setup();
});
