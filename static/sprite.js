function ZombieSprite() {
    this.state = 0;
    this.statesCount = 8;
    this.lastStateUpdate = 0;
    this.animetionTime = 100;
    this.scale = 0.5;
    this.size = 128;

    this.img = new Image();
    this.img.src = "static/zombie_topdown.png";
}

ZombieSprite.prototype.updateState = function (timeDiff) {
    this.lastStateUpdate += timeDiff;
    if (this.lastStateUpdate > this.animetionTime) {
        this.state = (this.state + 1) % this.statesCount;
        this.lastStateUpdate = 0;
    }
};

ZombieSprite.prototype.getMoveImage = function () {
    //context.drawImage(img, sx, sy, sw, sh, dx, dy, dw, dh)
    return [this.img, (4 + this.state)*this.size, this.size, this.size, this.size,
        -this.size*this.scale/2, -this.size*this.scale/2, this.size*this.scale, this.size*this.scale];
};