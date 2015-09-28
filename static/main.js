(function () {
    var url = "ws://" + window.location.hostname + ":7101/ws";
    var eventsHub = new EventsHub();
    var transport = new WebSocketTransport(url, eventsHub);
    var map = new PixiMap(eventsHub);
    var game = new Game(eventsHub, map);
    game.start();
})();