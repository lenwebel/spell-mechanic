const wordListElement = document.getElementById('word-list');
const socket = new WebSocket('ws://go-web-app:8765' + window.location.search);
const words = [];

console.log('Attempting Connection...');

socket.addEventListener('open', () => {
    console.log('Connected to WebSocket server');
});

socket.addEventListener('message', (event) => {
    const message = event.data; 
    console.log('Received message:', message);
    words.push(message);

    const c = wordListElement.value.split('\n').map(f => f.toLowerCase()).includes(message.toLowerCase()) ? 'valid' : 'invalid';

    document.getElementById('output').innerHTML += `<span class="${c}">` + message + '</span>';
});

socket.addEventListener('close', () => {
    console.log('Disconnected from WebSocket server');
});

function sendMsg() {
    const message = document.getElementById('message').value;
    socket.send(message);
    document.getElementById('message').value = '';
    message.setfocus();
}



window.onblur = () => {
    socket.send('Moved away from page')
}

window.onfocus =  () => {
    socket.send('Back on page')
}