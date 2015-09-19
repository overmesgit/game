(function () {
    var url = "ws://" + window.location.hostname + ":8000/ws";
    var eventsHub = new EventsHub();
    var transport = new WebSocketTransport(url, eventsHub);
    var map = new CanvasMap(eventsHub);
    var game = new Game(eventsHub, map);
    game.start();
})();