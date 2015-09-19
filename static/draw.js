function CanvasMap(eventsHub) {
    this.eventsHub = eventsHub;
    this.canvas = document.getElementById("map");
    if (this.canvas.getContext) {
        this.ctx = this.canvas.getContext("2d");
    }

    this.units = [];
    this.elemLeft = this.canvas.offsetLeft;
    this.elemTop = this.canvas.offsetTop;
    this.canvas.addEventListener('click', this.onClick.bind(this));
}

CanvasMap.prototype.onClick = function (event) {
    var x = event.pageX - this.elemLeft,
        y = event.pageY - this.elemTop;
    this.eventsHub.trigger(new Event('map:click', {'x': x, 'y': y}));
};

CanvasMap.prototype.unitsUpdate = function (units) {
    this.units = units;
    this.drawAllUnits();
    this.printUnitsCount();
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
    var rad = 2.5;
    switch (unit.Type) {
        case "fr":
            color = 'red';
            rad = 2;
            break
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
};

CanvasMap.prototype.printUnitsCount = function () {
    this.ctx.font = "30px Arial";
    this.ctx.fillText("Units: " + this.units.length, 10, 30);
};