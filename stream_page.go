package main

const streamPageHTML = `<html>
<!DOCTYPE html>
<html>
	<head>
	<script>
	// namespace MJPEG { ...
		var MJPEG = (function(module) {
			"use strict";
		
			// class Stream { ...
			module.Stream = function(args) {
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
						self.frameTimer = setInterval(function() {
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
		
				self.start = function() { setRunning(true); }
				self.stop = function() { setRunning(false); }
			};
		
			// class Player { ...
			module.Player = function(canvas, url, options) {
		
				var self = this;
				if (typeof canvas === "string" || canvas instanceof String) {
					canvas = document.getElementById(canvas);
				}
				var context = canvas.getContext("2d");
		
				if (! options) {
					options = {};
				}
				options.url = url;
				options.onFrame = updateFrame;
				options.onStart = function() { console.log("started"); }
				options.onStop = function() { console.log("stopped"); }
		
				self.stream = new module.Stream(options);
		
				canvas.addEventListener("click", function(e) {
					var xmlHttp = new XMLHttpRequest();
					xmlHttp.open( "GET", "http://localhost:8080/click/"+e.pageX+"/"+e.pageY, false ); 
					xmlHttp.send( null );
					return xmlHttp.responseText;
				}, false);
		
				function scaleRect(srcSize, dstSize) {
					var ratio = Math.min(dstSize.width / srcSize.width,
															 dstSize.height / srcSize.height);
					var newRect = {
						x: 0, y: 0,
						width: srcSize.width * ratio,
						height: srcSize.height * ratio
					};
					newRect.x = (dstSize.width/2) - (newRect.width/2);
					newRect.y = (dstSize.height/2) - (newRect.height/2);

					return newRect;
				}
		
				function updateFrame(img) {
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
						context.canvas.width  = window.innerWidth;
					  context.canvas.height = window.innerHeight;
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
						self.stop();
						console.log("!");
						throw e;
					}
				}
		
				self.start = function() { self.stream.start(); }
				self.stop = function() { self.stream.stop(); }
			};
		
			return module;
		
		})(MJPEG || {});
	</script>
	<meta charset="utf-8"/>
	<title>Player</title>
  </head>
  <body style="margin:0px;width:  100%; height: 100%;">
	<canvas id="player" style="background: #000; width: 100%; height:100%;">
	  Your browser sucks.
	</canvas>
  </body>
  
  <script>
	var player = new MJPEG.Player("player", "http://localhost:4545");
	player.start();
  </script>
</html>`
