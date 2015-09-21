function Unit(id, x, y, radius, type, h, speedX, speedY) {
    this.id = id;
    this.x = x;
    this.y = y;
    this.type = type;
    this.radius = radius;
    this.health = h;
    this.speedX = speedX;
    this.speedY = speedY;
    //this.directionX = directionX;
    //this.directionY = directionY;

    this.state = 0;
    this.statesCount = 8;
    this.lastStateUpdate = 0;
    this.animetionTime = 100;

    this.img = new Image();
    this.img.src = "static/zombie_topdown.png";

    this.types = {
        "Enemy": "en",
        "Bullet": "bu",
        "Player": "pl"
    };
}

Unit.prototype.update = function (x, y, h, speedX, speedY) {
    this.x = x;
    this.y = y;
    this.health = h;
    this.speedX = speedX;
    this.speedY = speedY;
};

Unit.prototype.updateState = function (timeDiff) {
    this.lastStateUpdate += timeDiff;
    if (this.lastStateUpdate > this.animetionTime) {
        this.state = (this.state + 1) % this.statesCount;
        this.lastStateUpdate = 0;
    }
};

Unit.prototype.getUnitColor = function () {
    var color = 'green';
    switch (this.type) {
        case this.types['Bullet']:
            color = 'red';
            break;
        case this.types['Player']:
            color = this.getPlayerColor();
            break;
    }
    if (this.health <= 0) {
        color = 'black';
    }
    return color;
};

Unit.prototype.getPlayerColor = function () {
    var color;
    switch (this.state) {
        case 0:
            color = 'deepskyblue';
            break;
        case 1:
            color = 'gold';
            break;
        case 2:
            color = 'olivedrab';
            break;
    }
    return color;
};

Unit.prototype.getImage = function () {
    //context.drawImage(img, sx, sy, sw, sh, dx, dy, dw, dh)
    return [this.img, (4 + this.state)*128, 0, 128, 128, this.x - 32, this.y - 32, 64, 64];
};