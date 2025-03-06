import './style.css'
import './app.css'

import { GetDirPath } from '../wailsjs/go/main/App'
import { GetAlbums } from '../wailsjs/go/main/App'

document.querySelector('#app').innerHTML = `
	<div class="input-box" id="player">
		<button class="btn" onclick="selectDirectory()">Select</button>
 </div>
`

window.selectDirectory = async function () {
	try {
		const dirPath = await GetDirPath()
		if (!dirPath) {
			console.log("User canceled directory selection");
			return;
		}
		console.log("Selected directory:", dirPath);
		try {
			let albums = GetAlbums(dirPath)
			displayAlbums(albums);
		} catch(err) {
			console.error("Error importing directory:", err);
		}
	} catch (err) {
		console.error("Error selecting directory:", err);
	}
}

function displayAlbums(albums) {
	console.log(albums)
}