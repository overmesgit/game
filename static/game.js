function Game(eventHub, map) {
    this.map = map;
    this.eventsHub = eventHub;
    this.state = 'ready';
}

Game.prototype.start = function() {
    this.eventsHub.on('ws:received', this.onWsMessage, this);
    this.eventsHub.on('ws:ready', this.onWsReady, this);
    this.eventsHub.on('map:click', this.onMapClick, this);

};

Game.prototype.onWsReady = function () {
    clearInterval(this.interval);
    this.state = 'ready';
    this.interval = setInterval(this.sendUnitsGet.bind(this), 50);
};

Game.prototype.sendUnitsGet = function () {
    this.state = 'wait';
    this.eventsHub.trigger(new Event('ws:send', {'get': 'units'}));
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

Game.prototype.onMapClick = function(event) {
    this.eventsHub.trigger(new Event('ws:send', {'get': 'boom', 'args': event.data}));
};