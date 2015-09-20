function CanvasMap(eventsHub) {
    this.eventsHub = eventsHub;
    this.canvas = document.getElementById("map");
    if (this.canvas.getContext) {
        this.ctx = this.canvas.getContext("2d");
    }

    this.mouseState = null;
    this.lastRenderTime = null;

    this.units = [];
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
    this.units = units;
    this.draw();
};

CanvasMap.prototype.draw = function () {
    this.drawAllUnits();
    this.printUnitsCount();
};

CanvasMap.prototype.unitsAnimate = function () {
    var timeDiff = (new Date().getTime() - this.lastRenderTime);
    for (var unit in this.units) {
        this.moveUnit(this.units[unit], timeDiff);
    }
    this.draw();
    window.requestAnimationFrame(this.unitsAnimate.bind(this));
};

CanvasMap.prototype.moveUnit = function (unit, step) {
    unit.X += unit.SX*step/1000;
    unit.Y += unit.SY*step/1000;
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

CanvasMap.prototype.painUnit = function (unit) {
    var color = 'green';
    var rad = unit.R;
    switch (unit.T) {
        case "fr":
            color = 'red';
            break;
        case "pl":
            color = 'blue';
            break;
    }
    this.painCircle(unit["X"], unit["Y"], rad, color);
};

CanvasMap.prototype.clear = function () {
    this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
};

CanvasMap.prototype.drawAllUnits = function () {
    this.clear();
    for (var unit in this.units) {
        this.painUnit(this.units[unit]);
    }
    this.lastRenderTime = new Date().getTime();
};

CanvasMap.prototype.printUnitsCount = function () {
    this.ctx.font = "30px Arial";
    this.ctx.fillText("Units: " + this.units.length, 10, 30);
};