import './style.css';
import './app.css';

import { Play } from '../wailsjs/go/main/App';

document.querySelector('#app').innerHTML = `
	<div class="input-box" id="player">
		Enter name: <input type="text" id="name"/>
        <button class="btn" onclick="play()">Play</button>
	</div>
`;

window.play = function() {
	let name = document.getElementById("name")["value"]
	console.log(name)
	try {
		Play(name);
	} catch (err) {
		console.error(err);
	}
}
