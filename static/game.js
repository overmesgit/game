function Game(eventHub, map) {
    this.map = map;
    this.eventsHub = eventHub;

    var body = document.getElementsByTagName('body')[0];
    body.onkeydown = this.keyDown.bind(this);
    body.onkeyup = this.keyUp.bind(this);
    this.state = 'ready';
    this.pressedKeys = {};
}

Game.prototype.start = function() {
    this.eventsHub.on('ws:received', this.onWsMessage, this);
    this.eventsHub.on('ws:ready', this.onWsReady, this);
    this.eventsHub.on('map:mousedown', this.onMouseDown, this);
    this.eventsHub.on('map:mouseup', this.onMouseUp, this);
    this.eventsHub.on('map:mousemove', this.onMouseMove, this);
    this.eventsHub.on('map:contextmenu', this.onMouseContext, this);
};

Game.prototype.onWsReady = function () {
    clearInterval(this.interval);
    this.state = 'ready';
    this.interval = setInterval(this.sendUnitsGet.bind(this), 50);
};

Game.prototype.sendUnitsGet = function () {
    this.state = 'wait';
    this.eventsHub.trigger('ws:send', {'get': 'units'});
};

Game.prototype.onWsMessage = function(event) {
    switch(event.data['get']) {
        case 'units':
            this.state = 'render';
            this.map.unitsUpdate(event.data['units']);
            break;
    }
    this.state = 'ready';
};

Game.prototype.onMouseDown = function(event) {
    this.eventsHub.trigger('ws:send', {'get': 'direction', 'args': event.data});
    this.eventsHub.trigger('ws:send', {'get': 'fire', 'args': true});
};

Game.prototype.onMouseUp = function(event) {
    this.eventsHub.trigger('ws:send', {'get': 'fire', 'args': false});
};

Game.prototype.onMouseContext = function(event) {
    this.eventsHub.trigger('ws:send', {'get': 'boom', 'args': event.data});
};

Game.prototype.onMouseMove = function(event) {
    if (this.map.mouseState == 'fire') {
        this.eventsHub.trigger('ws:send', {'get': 'direction', 'args': event.data});
    }
};

Game.prototype.keyDown = function(event) {
    var char = getChar(event);
    if(!event.repeat && 'WASD'.indexOf(char) >= 0) {
        this.pressedKeys[char] = true;
        this.eventsHub.trigger('ws:send', {'get': 'move', 'args': this.pressedKeys});
    }
};

Game.prototype.keyUp = function(event) {
    var char = getChar(event);
    if(!event.repeat && 'WASD'.indexOf(char) >= 0) {
        delete this.pressedKeys[char];
        this.eventsHub.trigger('ws:send', {'get': 'move', 'args': this.pressedKeys});
    }
};