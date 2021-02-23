function initWS() {
  var socket = new WebSocket("/socket")

  socket.onopen = function (e) {
    socket.send("Here is some text!")
  }

  socket.onmessage = function (e) {
    console.log(e.data);
  }
}

function openSocket() {
  console.log("opening socket!")
}

// Form submit event listener
const form = document.querySelector('form');
form.addEventListener('submit', e => {
  // Disable default action
  e.preventDefault();

  // Collect files
  const files = document.querySelector('[name=img]').files;

  // Set before image
  document.querySelector('#beforeImage').src = URL.createObjectURL(files[0]);

  // Get the socket url
  const socketURL = "ws://" + window.location.hostname + "/socket";

  // Open a websocket connection to the frontend socket endpoint
  var socket = new WebSocket(socketURL);

  socket.onopen = function (e) {
    socket.send(files[0]);
  }

  socket.onmessage = function (e) {
    var imageURL = URL.createObjectURL(e.data);
    document.querySelector('#responseImage').src = imageURL;
    socket.close();
  }
})
