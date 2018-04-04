package main

const streamPageHTML = `
<html>
<!DOCTYPE html>
<html>

<head>
    <meta charset="utf-8" />
    <title>{{appName}}</title>
    <script>
        var mouse = {
            x: 0,
            y: 0,
            leftButton: "none" //none, pressed, released, clicked
        }


        function sendMouseState() {
            // var xmlHttp = new XMLHttpRequest();
            // xmlHttp.open( "GET", "http://localhost:8080/click/"+mouse.x+"/"+mouse.y+"/"+mouse.leftButton, true ); 
            // xmlHttp.send( null );
            // return xmlHttp.responseText;
        }


        var MJPEG = (function (module) {
            "use strict";

            // class Stream { ...
            module.Stream = function (args) {
                var self = this;
                var autoStart = args.autoStart || false;

                self.url = args.url;
                self.refreshRate = args.refreshRate || 500;
                self.onStart = args.onStart || null;
                self.onFrame = args.onFrame || null;
                self.onStop = args.onStop || null;
                self.callbacks = {};
                self.running = false;
                self.frameTimer = 0;

                self.img = new Image();
                if (autoStart) {
                    self.img.onload = self.start;
                }
                self.img.src = self.url;

                function setRunning(running) {
                    self.running = running;
                    if (self.running) {
                        self.img.src = self.url;
                        self.frameTimer = setInterval(function () {
                            if (self.onFrame) {
                                self.onFrame(self.img);
                            }
                        }, self.refreshRate);
                        if (self.onStart) {
                            self.onStart();
                        }
                    } else {
                        self.img.src = "#";
                        clearInterval(self.frameTimer);
                        if (self.onStop) {
                            self.onStop();
                        }
                    }
                }

                self.start = function () { setRunning(true); }
                self.stop = function () { setRunning(false); }
            };

            // class Player { ...
            module.Player = function (canvas, url, options) {

                var self = this;
                if (typeof canvas === "string" || canvas instanceof String) {
                    canvas = document.getElementById(canvas);
                }
                var context = canvas.getContext("2d");

                if (!options) {
                    options = {};
                }
                options.url = url;
                options.onFrame = updateFrame;
                options.onStart = function () { console.log("started"); }
                options.onStop = function () { console.log("stopped"); }

                self.stream = new module.Stream(options);

                canvas.addEventListener("click", function (e) {
                    mouse.x = e.pageX
                    mouse.y = e.pageY

                    sendMouseState();
                }, false);

                function scaleRect(srcSize, dstSize) {
                    var ratio = Math.min(dstSize.width / srcSize.width,
                        dstSize.height / srcSize.height);
                    var newRect = {
                        x: 0, y: 0,
                        width: srcSize.width * ratio,
                        height: srcSize.height * ratio
                    };
                    newRect.x = (dstSize.width / 2) - (newRect.width / 2);
                    newRect.y = (dstSize.height / 2) - (newRect.height / 2);

                    return newRect;
                }

                function updateFrame(img) {
                    context.canvas.width = window.innerWidth;
                    context.canvas.height = window.innerHeight;
                    var srcRect = {
                        x: 0, y: 0,
                        width: img.naturalWidth,
                        height: img.naturalHeight
                    };
                    var dstRect = scaleRect(srcRect, {
                        width: canvas.width,
                        height: canvas.height
                    });
                    try {
                        context.drawImage(img,
                            0,
                            0,
                            srcRect.width,
                            srcRect.height,
                            0,
                            0,
                            context.canvas.width,
                            context.canvas.height
                        );
                        // console.log(".");

                    } catch (e) {
                        // if we can't draw, don't bother updating anymore
                        // self.stop();
                        console.log(e);
                        // throw e;
                    }
                }

                self.start = function () { self.stream.start(); }
                self.stop = function () { self.stream.stop(); }
            };

            return module;

        })(MJPEG || {});
    </script>
</head>

<body id="body" style="margin:0px;width:  100%; height: 100%;">
    <canvas id="player" style="background: #000; width: 100%; height:100%;">
        Your browser sucks.
    </canvas>
    <h1 id="ip" style="opacity:0.5;z-index:10;color:white;text-align: center; vertical-align:center;width: 100%;height: 100%;position: absolute;top: 0;left: 0;"></h1>

    <script>
        let parts = window.location.href.toString().split("/");
        document.title = parts[parts.length - 1];
    </script>
</body>

<script>
    var player = new MJPEG.Player("player", "http://localhost:4545", { refreshRate: 20 });
    player.start();
</script>

</html>
`
