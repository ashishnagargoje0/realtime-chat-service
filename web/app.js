let socket = new WebSocket(`ws://${window.location.host}/ws`);

socket.onmessage = function(event) {
    let li = document.createElement("li");
    li.textContent = event.data;
    document.getElementById("messages").appendChild(li);
};

function sendMessage() {
    let input = document.getElementById("msg");
    socket.send(input.value);
    input.value = "";
}
