function PixiMap(eventsHub) {
    this.eventsHub = eventsHub;

    this.units = {};

    var renderer = PIXI.autoDetectRenderer(1000, 800,{backgroundColor : 0x1099bb});

    document.getElementById("map").appendChild(renderer.view);
    var zombieImage = PIXI.Texture.fromImage('static/zombie_topdown.png');

    this.texture = new PIXI.Texture(zombieImage, new PIXI.Rectangle(0, 0, 128, 128));
    this.stage = new PIXI.Container();
    var stage = this.stage;

    this.basicText = new PIXI.Text(0);
    this.basicText.x = 20;
    this.basicText.y = 20;
    this.stage.addChild(this.basicText);

    this.stage.interactive = true;
    this.stage.on('mousedown', this.onMouseEvent, this);
    this.stage.on('mouseup', this.onMouseEvent, this);
    this.stage.on('mousemove', this.onMouseEvent, this);
    this.stage.on('rightclick', this.onMouseEvent, this);

    var lastStateUpdate = new Date().getTime();
    var map = this;

    animate();
    function animate() {
        requestAnimationFrame(animate);

        var timeDiff = (new Date().getTime() - lastStateUpdate);
        map.moveUnits(timeDiff);
        renderer.render(stage);
        lastStateUpdate = new Date().getTime();

    }
}

PixiMap.prototype.moveUnits = function (step) {
    for (var i in this.units) {
        var unit = this.units[i];
        unit.position.x += unit.speedX*step/1000;
        unit.position.y += unit.speedY*step/1000;
    }
};

PixiMap.prototype.onMouseEvent = function (event) {
    if (event.type == 'mousedown') this.mouseState = 'fire';
    if (event.type == 'mouseup') this.mouseState = null;

    var x = event.data.global.x,
        y = event.data.global.y;
    this.eventsHub.trigger('map:' + event.type, {'x': x, 'y': y});
};

PixiMap.prototype.unitsUpdate = function (units) {
    var newZombies = {};
    units.forEach(function(u){
        newZombies[u.ID] = true;
        var zombie = null;
        if (u.ID in this.units) {
            zombie = this.units[u.ID];
            if (u.H <= 0) {
                this.stage.removeChild(zombie);
                delete this.units[u.ID];
            }
        } else {
            zombie = new PIXI.Sprite(this.texture);
            zombie.anchor.x = 0.5;
            zombie.anchor.y = 0.5;
            this.stage.addChild(zombie);
            this.units[u.ID] = zombie;
        }
        zombie.position.x = u.X;
        zombie.position.y = u.Y;
        zombie.speedX = u.SX;
        zombie.speedY = u.SY;
    }.bind(this));

    for (var i in this.units) {
        if (!(i in newZombies)) {
            this.stage.removeChild(this.units[i]);
            delete this.units[i];
        }
    }

    this.basicText.text = units.length;
};