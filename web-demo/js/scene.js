let createScene = function() {

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

    let scene = null;
    let camera = null;
    let renderer = null;
    let objects = {}

    return {
        build: function(data, displayWidth, displayHeight)
        {
            scene = new THREE.Scene();
            camera = this.createCamera();

            for (let i = 0; i < data.length; i++) {
                let datum = data[i];
                let object = this.createObject(datum);
                objects[datum.E] = object
                scene.add(object);
            }

            renderer = new THREE.WebGLRenderer({ antialias: true });
            renderer.setSize(displayWidth, displayHeight);

            monitor.innerHTML = "";
            monitor.appendChild( renderer.domElement );

            renderer.render( scene, camera );
        },

        update: function(data) {
            for (let i = 0; i < data.length; i++) {
                let datum = data[i];
                let x = (datum.X / scale);
                let y = (datum.Y / scale);
                let z = -(datum.Z / scale);
                let object = objects[datum.E]
                object.position.set(x, y, z);
            }

            renderer.render( scene, camera );
        },

        // https://stackoverflow.com/questions/26021618/how-can-i-create-an-axonometric-oblique-cavalier-cabinet-with-threejs
        createCamera() {
            var camera = new THREE.OrthographicCamera( -10, 10, 5, -5, 0, 1000 );
            var matrix = new THREE.Matrix4();

            var alpha = Math.PI / 6;

            var Syx = 0,
                Szx = - 0.5 * Math.cos( alpha ),
                Sxy = 0,
                Szy = - 0.5 * Math.sin( alpha ),
                Sxz = 0,
                Syz = 0;

            matrix.set(1, Syx, Szx,0, Sxy,1, Szy,0, Sxz, Syz,1,0,0,0,0,1);

            camera.projectionMatrix.multiply(matrix);
            camera.projectionMatrixInverse.getInverse( camera.projectionMatrix );
            camera.position.set(8, 4.5, 0);

            return camera;
        },

        // https://stackoverflow.com/questions/36472653/drawing-edges-of-a-mesh-in-three-js
        createEdges: function(mesh)
        {
            var geometry = new THREE.EdgesGeometry( mesh.geometry );
            var material = new THREE.LineBasicMaterial( { color: 0xf0f0f0 } );

            return new THREE.LineSegments( geometry, material );
        },

        createObject: function(datum)
        {
            if (datum.Type === "handXXX") {
                return this.createHand(datum)
            } else if (datum.Type === "pyramid") {
                return this.createPyramid(datum)
            } else {
                return this.createBlock(datum)
            }
        },

        createHand: function(datum)
        {
            const group = new THREE.Group();

            const geometry = new THREE.BoxGeometry( 1, 1, 1 );
            const material = new THREE.MeshBasicMaterial( {color: 0x00ff00} );
            const cube = new THREE.Mesh( geometry, material );
            group.add( cube );

            return group;
        },

        createPyramid: function(datum)
        {
            var geometry = new THREE.Geometry();

            let x = (datum.X / scale);
            let y = (datum.Y / scale);
            let z = -(datum.Z / scale);
            let w = (datum.Width / scale);
            let l = -(datum.Length / scale);
            let h = (datum.Height / scale);

            geometry.vertices = [
                new THREE.Vector3(  0, 0, 0 ),
                new THREE.Vector3( w, 0, 0 ),
                new THREE.Vector3( w, 0, l ),
                new THREE.Vector3( 0, 0, l ),
                new THREE.Vector3( 0.5 * w, h, 0.5 * l )
            ];

            geometry.faces = [
                new THREE.Face3( 0, 1, 2 ),
                new THREE.Face3( 0, 2, 3 ),
                new THREE.Face3( 1, 0, 4 ),
                new THREE.Face3( 2, 1, 4 ),
                new THREE.Face3( 3, 2, 4 ),
                new THREE.Face3( 0, 3, 4 )
            ];

            let material = new THREE.MeshBasicMaterial( {color: colors[datum.Color] , wireframe:false, transparent: true, opacity: opacity} );
            let object = new THREE.Mesh( geometry, material );

            let group = new THREE.Group();

            group.add(object);

            let edges = this.createEdges(object)
            group.add(edges)

            group.position.set(x, y, z);

            return group;
        },

        createBlock: function(datum)
        {
            var geometry = new THREE.Geometry();

            let x = (datum.X / scale);
            let y = (datum.Y / scale);
            let z = -(datum.Z / scale);
            let w = (datum.Width / scale);
            let l = -(datum.Length / scale);
            let h = ((datum.Height ? datum.Height : 0.01) / scale);

            geometry.vertices = [
                new THREE.Vector3( 0, 0, 0 ),
                new THREE.Vector3( w, 0, 0 ),
                new THREE.Vector3( w, 0, l ),
                new THREE.Vector3( 0, 0, l ),

                new THREE.Vector3( 0, h, 0 ),
                new THREE.Vector3( w, h, 0 ),
                new THREE.Vector3( w, h, l ),
                new THREE.Vector3( 0, h, l ),
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
            let material = new THREE.MeshBasicMaterial( {color: colors[datum.Color] , wireframe:false, transparent: true, opacity: blockOpacity} );
            let object = new THREE.Mesh( geometry, material );
            let edges = this.createEdges(object)

            let group = new THREE.Group();
            group.add(object);
            group.add(edges)

            group.position.set(x, y, z);

            return group;
        }
    }
};