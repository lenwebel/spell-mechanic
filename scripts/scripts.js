const wordListElement = document.getElementById('word-list');
const spellout = document.getElementById('spellout');
const socket = new WebSocket('ws://192.168.86.22:8765' + window.location.search);
const words = [];

console.log('Attempting Connection...');

socket.addEventListener('open', () => {
    console.log('Connected to WebSocket server');
});

socket.addEventListener('message', (event) => {
    const message = event.data;
    console.log('Received message:', message);
    words.push(message);
    
    const c = wordListElement?.value?.split('\n').map(f => f.toLowerCase()).includes(message.toLowerCase()) ? 'valid' : 'invalid';

    const out = document.getElementById('output')
    if(out)
        out.innerHTML += `<span class="${c}">` + message ?? '' + '</span>';
});

socket.addEventListener('close', () => {
    console.log('Disconnected from WebSocket server');
});

function sendMsg() {
    const message = document.getElementById('message').value;

    if(!message) return;
    socket.send(message);
    document.getElementById('message').value = '';
}

window.onblur = () => {
    socket.send('Moved away from page')
}

window.onfocus = () => {
    socket.send('Back on page')
}


function onMessageChanged(e) {
    if (spellout)
        spellout.innerHTML = e.target.value;
}

message.removeEventListener('keyup',(e) => onMessageChanged(e));
message.addEventListener('keyup', (e) => onMessageChanged(e));
