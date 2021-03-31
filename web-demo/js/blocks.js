$(function(){

    const inputField = document.getElementById('q');
    const samplePopup = document.getElementById('samples');
    const answerBox = document.getElementById('answer-box');
    const errorBox = document.getElementById('error-box');
    const optionPopupCloseButton = document.getElementById('close');
    const sampleButton = document.getElementById('show-samples');
    const form = document.getElementById('f');
    const optionsPopup = document.getElementById('options-popup');
    const optionsHeader = document.getElementById('options-header');
    const optionsHtml = document.getElementById('options-box');
    const logBox = document.getElementById("log");
    const monitor = document.getElementById("monitor");
    const resetButton = document.getElementById("reset");

    const displayWidth = Math.min(window.innerWidth, 600);
    const displayHeight = displayWidth / 2;

    let scene = createScene();
    let currentInput = "";

    function setup()
    {
        monitor.style.height = displayHeight + "px";

        optionPopupCloseButton.onclick = function() {
            hideOptionsPopup()
        };

        sampleButton.onclick = function () {
            showSamplePopup()
        };

        resetButton.onclick = function () {
            reset();
        }

        form.onsubmit = function(){
            currentInput = inputField.value
            clearProductions()
            clearInput();
            tell(currentInput)
            return false;
        };

        document.addEventListener('keydown', function(event){
            if (event.key === 'Escape') {
                hideSamplePopup()
                hideOptionsPopup()
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

    function clearProductions()
    {
        document.getElementById('Intro').innerHTML = "";
        document.getElementById('Ready').innerHTML = "";
        document.getElementById('Tokens').innerHTML = "";
        document.getElementById('Parse-tree').innerHTML = "";
        document.getElementById('Relations').innerHTML = "";
        document.getElementById('Solution').innerHTML = "";
        document.getElementById('Answer').innerHTML = "";
        document.getElementById('Anaphora-queue').innerHTML = "";
    }

    function showSamplePopup() {
        samplePopup.style.display = "block";
    }

    function hideSamplePopup() {
        samplePopup.style.display = "none";
    }

    function errorToHtml(error) {
        let html = "";

        if (Array.isArray(error)) {
            for (let i = 0; i < error.length; i++) {
                html += error[i] + "<br>";
            }
        } else {
            html = error;
        }
        return html;
    }

    function showError(error) {
        errorBox.innerHTML = errorToHtml(error);
    }

    function showAnswer(answer) {
        answerBox.innerHTML = answer;
    }

    function clearInput() {
        inputField.value = "";
    }

    function showProductions(productions) {

        for (let key in productions) {
            let production = productions[key];

            let matches = production.match(/([^:]+)/);
            let name = matches[1];
            let value = production.substr(name.length + 2);
            let id = name.replace(' ', '-')

            let container = document.getElementById(id)
            if (container) {
                container.innerHTML = "<div class='card'><h2>" + name + "</h2>" + "<pre>" + value + "</pre></div>";
            }
        }
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
            type: 'POST',
            success: function (data) {

                if (data.ErrorLines.length > 0) {
                    showAnswer("")
                    showError(data.ErrorLines);
                    log(currentInput, errorToHtml(data.ErrorLines))
                } else {
                    processResponse(data.Message)
                    showError([]);
                }
                showProductions(data.Productions);
            },
            error: function (request, status, error) {
                showError(error)
            }
        });
    }

    function processResponse(response)
    {
        let asserts = [];
        let assert;
        let maxDuration = 0;
        let animations = [];

        for (let i = 0; i < response.length; i++) {
            let relation = response[i];
            switch (relation.predicate) {
                case 'dom_action_move_to':
                    assert = moveObject(relation)
                    asserts.push(assert)
                    let result = scene.createObjectAnimation({
                        E: "`" + relation.arguments[1].sort + ":" + relation.arguments[1].value + "`",
                        X: relation.arguments[2].value,
                        Y: relation.arguments[4].value,
                        Z: relation.arguments[3].value
                    })
                    animations.push(result.animation)
                    maxDuration = Math.max(maxDuration, result.duration)
                    break;
                case 'go_print':
                    assert = print(relation)
                    asserts.push(assert)
                    break;
                case 'go_user_select':
                    showOptionsPopup(relation)
                    break;
            }
        }

        if (animations.length > 0) {
            scene.runAnimations(animations, maxDuration)
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
        log(currentInput, answer)
        return getAssert(relation);
    }

    function moveObject(relation) {
        return getAssert(relation);
    }

    function tell(input) {
        sendRequest([{
            predicate: 'go_tell',
            arguments: [
                {
                    type: 'string',
                    value: input
                }
            ]
        }]);
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

    function showOptionsPopup(relation) {

        let options = relation.arguments[0].list;

        let html = "<ol>";
        for (let i = 0; i < options.length; i++) {
            let argument = options[i];
            html += "<li><a href='" + i + "'>" + argument.value + "</a></li>";
        }
        html += "</ol>"

        optionsHtml.innerHTML = html;
        optionsPopup.style.display = "block";

        let aTags = optionsHtml.querySelectorAll('a');
        for (let i = 0; i < aTags.length; i++) {
            aTags[i].onclick = function (event) {
                event.preventDefault();
                hideOptionsPopup()
                relation.arguments[1].type = "string"
                relation.arguments[1].value = event.currentTarget.getAttribute('href');
                let assert = getAssert(relation)
                sendRequest([assert])
            };
        }
    }

    function hideOptionsPopup()
    {
        optionsPopup.style.display = "none";
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
