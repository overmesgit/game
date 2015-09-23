function CanvasMap(eventsHub) {
    this.eventsHub = eventsHub;
    this.canvas = document.getElementById("map");
    if (this.canvas.getContext) {
        this.ctx = this.canvas.getContext("2d");
    }

    this.mouseState = null;
    this.lastStateUpdate = null;

    this.units = {};
    this.elemLeft = this.canvas.offsetLeft;
    this.elemTop = this.canvas.offsetTop;
    this.canvas.addEventListener('mousedown', this.onMouseEvent.bind(this));
    this.canvas.addEventListener('mouseup', this.onMouseEvent.bind(this));
    this.canvas.addEventListener('mousemove', this.onMouseEvent.bind(this));
    this.canvas.addEventListener('contextmenu', this.onMouseEvent.bind(this));

    window.requestAnimationFrame(this.unitsAnimate.bind(this));
}

CanvasMap.prototype.onMouseEvent = function (event) {
    if (event.type == 'mousedown') this.mouseState = 'fire';
    if (event.type == 'mouseup') this.mouseState = null;

    var x = event.pageX - this.elemLeft,
        y = event.pageY - this.elemTop;
    this.eventsHub.trigger('map:' + event.type, {'x': x, 'y': y});
};

CanvasMap.prototype.unitsUpdate = function (units) {
    units.forEach(function(u){
        if (u.ID in this.units) {
            this.units[u.ID].update(u.X, u.Y, u.H, u.SX, u.SY);
        } else {
            this.units[u.ID] = new Unit(u.id, u.X, u.Y, u.R, u.T, u.H, u.SX, u.SY)
        }
    }.bind(this));
};

CanvasMap.prototype.draw = function (timeDiff) {
    this.drawAllUnits(timeDiff);
    this.printUnitsCount();
};

CanvasMap.prototype.unitsAnimate = function () {
    var timeDiff = (new Date().getTime() - this.lastStateUpdate);
    for (var i in this.units) {
        this.moveUnit(this.units[i], timeDiff);
    }
    this.draw(timeDiff);
    window.requestAnimationFrame(this.unitsAnimate.bind(this));
};

CanvasMap.prototype.collapseUnit = function (unit, step) {
    unit.radius -= 0.003*step;
    if (unit.radius < 0) {
        unit.radius = 0;
    }
};

CanvasMap.prototype.moveUnit = function (unit, step) {
    unit.x += unit.speedX*step/1000;
    unit.y += unit.speedY*step/1000;
};

CanvasMap.prototype.painCircle = function (x, y, r, color) {
    this.ctx.beginPath();
    this.ctx.arc(x, y, r, 0, 2 * Math.PI, false);
    this.ctx.fillStyle = color;
    this.ctx.fill();
    this.ctx.lineWidth = 1;
    this.ctx.strokeStyle = '#003300';
    this.ctx.stroke();
};

CanvasMap.prototype.painUnit = function (unit, timeDiff) {
    unit.sprite.updateState(timeDiff);
    if (unit.type == 'en') {
        this.painEnemy(unit);
    } else {
        var color = unit.getUnitColor();
        var rad = unit.radius;
        this.painCircle(unit.x, unit.y, rad, color);
    }

};

CanvasMap.prototype.painEnemy = function (unit) {
    var unitImg = unit.sprite.getMoveImage();
    var angleRadians = Math.atan2(unit.y - (unit.y + unit.speedY), unit.x - (unit.x + unit.speedX));
    this.ctx.translate(unit.x, unit.y);
    this.ctx.rotate(angleRadians);
    //context.drawImage(img, sx, sy, sw, sh, dx, dy, dw, dh)
    this.ctx.drawImage.apply(this.ctx, unitImg);
    this.ctx.rotate(-angleRadians);
    this.ctx.translate(-unit.x, -unit.y);
};

CanvasMap.prototype.clear = function () {
    this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
};

CanvasMap.prototype.drawAllUnits = function (timeDiff) {
    this.clear();
    for(var i in this.units) {
        if (this.units[i].health > 0) {
            this.painUnit(this.units[i], timeDiff);
        } else {
            delete this.units[i];
        }
    }
    this.lastStateUpdate = new Date().getTime();
};

CanvasMap.prototype.printUnitsCount = function () {
    this.ctx.font = "30px Arial";
    this.ctx.fillText("Units: " + 0, 10, 30);
};