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

    let webSocket = null;

    function setup() {

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

        document.addEventListener('keydown', function(event) {
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

        webSocket = new WebSocket("ws://127.0.0.1:3333/");
        webSocket.onopen = () => {
            updateScene(true)
        }
        webSocket.onmessage = (event) => {
            handleIncomingMessage(JSON.parse(event.data))
        }
    }

    function handleIncomingMessage(message) {
        console.log(message)
        switch (message.MessageType) {
            case "description":
                scene.build(monitor, message.Message, displayWidth, displayHeight)
                break
            case "print":
                print(message.Message)
                send("language", "acknowledge", "")
                break
            case "move_to":
                doMoveTo(message.Resource, message.Message[0])
                break
        }
    }

    function doMoveTo(resource, moves) {
        let maxDuration = 0;
        let animations = [];

        for (const move of moves) {
            let result = scene.createObjectAnimation({
                E: move[0],
                X: move[1],
                Y: move[3],
                Z: move[2]
            })
            animations.push(result.animation)
            maxDuration = Math.max(maxDuration, result.duration)
        }

        if (animations.length > 0) {
            scene.runAnimations(animations, maxDuration)
        }
        window.setTimeout(function (){
            send("robot", "acknowledge", "")
        }, maxDuration);
    }

    function clearProductions() {
        document.getElementById('Intro').innerHTML = "";
        document.getElementById('Ready').innerHTML = "";
        document.getElementById('Tokens').innerHTML = "";
        document.getElementById('Parse-tree').innerHTML = "";
        document.getElementById('Relations').innerHTML = "";
        document.getElementById('Solution').innerHTML = "";
        document.getElementById('Answer').innerHTML = "";
        document.getElementById('Dialog-entities').innerHTML = "";
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
    }


    function print(answer) {
        showAnswer(answer)
        log(currentInput, answer)
    }

    function send(resource, messageType, message) {
        console.log("send", messageType, message)
        webSocket.send(JSON.stringify({
            System: "blocks",
            Resource: resource,
            MessageType: messageType,
            Message: message
        }))
    }

    function tell(input) {
        send("language", "respond", input)
    }

    function showOptionsPopup(relation) {

        let options = relation.arguments[1].list;

        let html = "<ol>";
        for (let i = 0; i < options.length; i++) {
            let argument = options[i];
            html += "<li><a href='" + i + "'>" + argument.value + "</a></li>";
        }
        html += "</ol>"

        optionsHeader.innerHTML = relation.arguments[0].value
        optionsHtml.innerHTML = html;
        optionsPopup.style.display = "block";

        let aTags = optionsHtml.querySelectorAll('a');
        for (let i = 0; i < aTags.length; i++) {
            aTags[i].onclick = function (event) {
                event.preventDefault();
                hideOptionsPopup()
                relation.arguments[2].type = "string"
                relation.arguments[2].value = event.currentTarget.getAttribute('href');
                let assert = getAssert(relation)
                sendRequest([assert])
            };
        }
    }

    function hideOptionsPopup() {
        optionsPopup.style.display = "none";
    }

    function log(question, answer) {
        let html = "";

        html += "<div class='prev-question'><h3>" + question + "</h3></div>";
        html += "<div class='prev-answer'>" + answer + "</div>";

        logBox.innerHTML = html + logBox.innerHTML;
    }

    function updateScene(initial) {
        send("no-resource", "describe", '')
    }

    setup();
});
