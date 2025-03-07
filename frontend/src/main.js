import './style.css'
import './app.css'

import { GetDirPath } from '../wailsjs/go/main/App'
import { GetAlbums } from '../wailsjs/go/main/App'
import { GetTrackInfo } from '../wailsjs/go/main/App'

document.querySelector('#app').innerHTML = `
	<div class="input-box" id="player">
		<button class="btn" onclick="trackInfo()">Select</button>
 </div>
`

window.trackInfo = async function () {
	try {
		result = GetTrackInfo("/Users/stefannosek/Documents/muzyka/Discovery (Daft Punk, 2001)/01. One More Time.flac")
		console.log(result)
	} catch (err) {
		console.error("Error fetching metadata: ", err);
	}
}

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