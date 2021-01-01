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
            popup.style.display = "none";
        };

        sampleButton.onclick = function (event) {
            samplePopup.style.display = "block";
        };

        resetButton.onclick = function () {
            reset();
        }

        form.onsubmit = function(){
            postQuestion(inputField.value);
            return false;
        };

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
            let value = production.substr(name.length + 1)
                .replace(/&/g, "&amp;")
                .replace(/</g, "&lt;")
                .replace(/>/g, "&gt;")
                .replace(/"/g, "&quot;")
                .replace(/'/g, "&#039;")
                .replace(/\n/g, "<br>");

            if (name === "Parse tree" || name === "Relations") {
                value = value.replace(/  /g, "&nbsp;&nbsp;")
            }

            html[container] += "<div class='card'><h2>" + name + "</h2>" + "<p>" + value + "</p></div>";

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

    function postQuestion(question) {
        $.ajax({
            url: 'ajax-answer.php',
            data: { format: "json", query: question, app: "blocks" },
            dataType: 'json',
            type: 'GET',
            success: function (data) {

                if (data.OptionKeys.length === 0) {
                    showAnswer(data.Answer);
                    clearInput();
                    log(question, data.Answer)
                    updateScene(false)
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

        optionsHeader.innerHTML = answer;
        optionsBox.innerHTML = html;

        popup.style.display = showOptions ? "block" : "none";

        let aTags = optionsBox.querySelectorAll('a');
        for (let i = 0; i < aTags.length; i++) {
            aTags[i].onclick = function (event) {
                event.preventDefault();
                postQuestion(event.currentTarget.getAttribute('href'));
            };
        }
    }

    function log(question, answer) {
        let html = "";

        html += "<div><h3>" + question + "</h3></div>";
        html += "<div>" + answer + "</div>";

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
                    scene.build(data, displayWidth, displayHeight)
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
