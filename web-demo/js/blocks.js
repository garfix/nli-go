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
    const displayWidth = Math.min(window.innerWidth, 600);
    const displayHeight = displayWidth / 2;
    const colors = {
        red: 0xc00000,
        green: 0x008000,
        blue: 0x0000c0,
        white: 0xc0c0c0,
        black: 0x654321
    }

    const scale = 90;
    const opacity = 0.9;
    const boxOpacity = 0.4;

    function setup()
    {
        monitor.style.height = displayHeight + "px";

        popupCloseButton.onclick = function() {
            popup.style.display = "none";
        };

        sampleButton.onclick = function (event) {
            event.preventDefault();
            samplePopup.style.display = "block";
        };

        form.onsubmit = function(){
            postQuestion(inputField.value);
            return false;
        };

        let samples = document.querySelectorAll('#samples li');
        for (let i = 0; i < samples.length; i++) {
            samples[i].onclick = function (element) {
                let li = element.currentTarget;
                inputField.value = li.innerHTML;
                samplePopup.style.display = "none";
            }
        }

        updateMonitor()
    }

    function showError(error) {
        let html = "";

        for (let i = 0; i < error.length; i++) {
            html += error[i] + "<br>";
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

            if (name === "Parse tree") {
                value = value.replace(/ /g, "&nbsp;")
            }

            html[container] += "<h2>" + name + "</h2>" + "<p>" + value + "</p>";

            if (name === "Parse tree") {
                container = 'production-box-right';
            }
        }

        productionBoxLeft.innerHTML = html['production-box-left'];
        productionBoxRight.innerHTML = html['production-box-right'];
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
                    updateMonitor()
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

    function updateMonitor()
    {
        $.ajax({
            url: 'scene.php',
            data: { format: "json" },
            dataType: 'json',
            type: 'GET',
            success: function (data) {
                buildScene(data)
            },
            error: function (request, status, error) {
                showError(error)
            }
        });
    }

    function buildScene(data)
    {
        const scene = new THREE.Scene();
        var camera = new THREE.OrthographicCamera( -10, 10, 5, -5, 0, 1000 );

        const group = new THREE.Group();

        for (let i = 0; i < data.length; i++) {
            let datum = data[i];
            let object = createObject(datum);
            skew(object)
            group.add(createEdges(object))
            group.add(object);
        }

        scene.add( group );

        const directionalLight = new THREE.DirectionalLight( 0xffffff, 0.5 );
        scene.add( directionalLight );

        camera.position.set(8, 4.5, 0);
        const renderer = new THREE.WebGLRenderer({ antialias: true });
        renderer.setSize(displayWidth, displayHeight);

        monitor.innerHTML = "";
        monitor.appendChild( renderer.domElement );

        renderer.render( scene, camera );
    }

    // https://stackoverflow.com/questions/36472653/drawing-edges-of-a-mesh-in-three-js
    function createEdges(mesh)
    {
        var geometry = new THREE.EdgesGeometry( mesh.geometry );
        var material = new THREE.LineBasicMaterial( { color: 0xf0f0f0 } );

        return new THREE.LineSegments( geometry, material );
    }

    // https://stackoverflow.com/questions/26021618/how-can-i-create-an-axonometric-oblique-cavalier-cabinet-with-threejs
    function skew(mesh)
    {
        var alpha = Math.PI / 6;

        var Syx = 0,
            Szx = - 0.5 * Math.cos( alpha ),
            Sxy = 0,
            Szy = - 0.5 * Math.sin( alpha ),
            Sxz = 0,
            Syz = 0;

        var matrix = new THREE.Matrix4();

        matrix.set(   1,   Syx,  Szx,  0,
            Sxy,     1,  Szy,  0,
            Sxz,   Syz,   1,   0,
            0,     0,   0,   1  );

        mesh.geometry.applyMatrix4(matrix);
    }

    function createObject(datum)
    {
        if (datum.Type === "pyramid") {
            return createPyramid(datum)
        } else {
            return createBlock(datum)
        }
    }

    function createPyramid(datum)
    {
        var geometry = new THREE.Geometry();

        let x = (datum.X / scale);
        let y = (datum.Y / scale);
        let z = -(datum.Z / scale);
        let w = (datum.Width / scale);
        let l = -(datum.Length / scale);
        let h = (datum.Height / scale);

        geometry.vertices = [
            new THREE.Vector3( x + 0, y + 0, z + 0 ),
            new THREE.Vector3( x + w, y + 0, z + 0 ),
            new THREE.Vector3( x + w, y + 0, z + l ),
            new THREE.Vector3( x + 0, y + 0, z + l ),
            new THREE.Vector3( x + 0.5 * w, y + h, z + 0.5 * l )
        ];

        geometry.faces = [
            new THREE.Face3( 0, 1, 2 ),
            new THREE.Face3( 0, 2, 3 ),
            new THREE.Face3( 1, 0, 4 ),
            new THREE.Face3( 2, 1, 4 ),
            new THREE.Face3( 3, 2, 4 ),
            new THREE.Face3( 0, 3, 4 )
        ];

        var material = new THREE.MeshBasicMaterial( {color: colors[datum.Color] , wireframe:false, transparent: true, opacity: opacity} );
        return new THREE.Mesh( geometry, material );
    }

    function createBlock(datum)
    {
        var geometry = new THREE.Geometry();

        let x = (datum.X / scale);
        let y = (datum.Y / scale);
        let z = -(datum.Z / scale);
        let w = (datum.Width / scale);
        let l = -(datum.Length / scale);
        let h = ((datum.Height ? datum.Height : 0.01) / scale);

        geometry.vertices = [
            new THREE.Vector3( x + 0, y + 0, z + 0 ),
            new THREE.Vector3( x + w, y + 0, z + 0 ),
            new THREE.Vector3( x + w, y + 0, z + l ),
            new THREE.Vector3( x + 0, y + 0, z + l ),

            new THREE.Vector3( x + 0, y + h, z + 0 ),
            new THREE.Vector3( x + w, y + h, z + 0 ),
            new THREE.Vector3( x + w, y + h, z + l ),
            new THREE.Vector3( x + 0, y + h, z + l ),
        ];

        geometry.faces = [
            new THREE.Face3(0 ,1,  2),
            new THREE.Face3(2, 3,  0),

            new THREE.Face3( 0, 1,  5),
            new THREE.Face3( 5, 4,  0),

            new THREE.Face3( 0, 3,  7),
            new THREE.Face3( 7, 4,  0),

            new THREE.Face3( 1, 2,  6),
            new THREE.Face3( 6, 5,  1),

            new THREE.Face3( 3, 2,  6),
            new THREE.Face3( 6, 7,  3),

            new THREE.Face3( 6, 5,  4),
            new THREE.Face3( 4, 7,  6),
        ];

        let blockOpacity = datum.Type === "box" ? boxOpacity : opacity;

        var material = new THREE.MeshBasicMaterial( {color: colors[datum.Color] , wireframe:false, transparent: true, opacity: blockOpacity} );
        return new THREE.Mesh( geometry, material );
    }

    setup();
});
