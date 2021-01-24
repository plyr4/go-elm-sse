import './main.css';
import { Elm } from './Main.elm';

var app = Elm.Main.init({
  node: document.getElementById('root')
});


app.ports.connectSSE.subscribe(function (config) {

  console.log("creating new EventSource client");

  // create a new client using configuration URL
  var stream = new EventSource(config.url);

  // add an event listener to subscribe to incoming events
  stream.addEventListener("message", function (e) {
    app.ports.messageReceiver.send(e.data);
  });
});
