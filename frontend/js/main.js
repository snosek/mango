import { GetCatalog, GetAlbum, GetDirPath } from '../wailsjs/go/main/App';
import { catalog } from '../wailsjs/go/models.ts';
import { renderAlbumsList, renderAlbumDetails } from './album.js';
import { pauseSong, playSong } from './player.js';

let state = {
	currentView: 'albums',
	currentAlbum: null,
	currentTrack: null,
	isPlaying: false,
	catalog: catalog.Catalog
};

async function init() {
	document.getElementById('select-dir-button').addEventListener('click', handleSelectDirectory);
	document.getElementById('back-button').addEventListener('click', navigateToAlbums);
	document.getElementById('albums-container').addEventListener('click', handleAlbumClick);
	document.getElementById('tracks-list').addEventListener('click', handleTrackClick);
	document.getElementById('play-button').addEventListener('click', handlePlayClick);
	document.getElementById('pause-button').addEventListener('click', handlePauseClick);

	await loadAlbums();
}

async function loadAlbums(fp) {
	try {
		state.catalog = await GetCatalog(fp);
		renderAlbumsList(state.catalog.Albums, document.getElementById('albums-container'));
		navigateToAlbums();
	} catch (error) {
		console.error('Failed to load albums:', error);
		alert('Failed to load albums. Please try again.');
	} 
}

async function handleSelectDirectory() {
	try {
		const dirPath = await GetDirPath();
		if (dirPath) {
			await loadAlbums(dirPath);
		}
	} catch (error) {
		console.error('Error selecting directory:', error);
		alert('Failed to select directory: ' + error);
	} 
}

async function handleAlbumClick(event) {
	const albumCard = event.target.closest('.album-card');
	if (!albumCard) 
		return;
	const albumId = albumCard.dataset.id;
	navigateToAlbumDetails(albumId);
}

function handlePlayClick() {
	console.log(memorySizeOf(state.catalog))
	playSong(state.currentAlbum.Tracks[0])
}

function memorySizeOf(obj) {
	var bytes = 0;

	function sizeOf(obj) {
		if (obj !== null && obj !== undefined) {
			switch (typeof obj) {
				case "number":
					bytes += 8;
					break;
				case "string":
					bytes += obj.length * 2;
					break;
				case "boolean":
					bytes += 4;
					break;
				case "object":
					var objClass = Object.prototype.toString.call(obj).slice(8, -1);
					if (objClass === "Object" || objClass === "Array") {
						for (var key in obj) {
							if (!obj.hasOwnProperty(key)) continue;
							sizeOf(obj[key]);
						}
					} else bytes += obj.toString().length * 2;
					break;
			}
		}
		return bytes;
	}

	function formatByteSize(bytes) {
		if (bytes < 1024) return bytes + " bytes";
		else if (bytes < 1048576) return (bytes / 1024).toFixed(3) + " KiB";
		else if (bytes < 1073741824) return (bytes / 1048576).toFixed(3) + " MiB";
		else return (bytes / 1073741824).toFixed(3) + " GiB";
	}

	return formatByteSize(sizeOf(obj));
}

function handlePauseClick() {
	pauseSong(state.currentAlbum.Tracks[0])
}

function handleTrackClick(event) {
	const trackItem = event.target.closest('.track-item');
	if (!trackItem) return;

	const trackIndex = parseInt(trackItem.dataset.index, 10);
	playTrack(state.currentAlbum.Tracks[trackIndex]);
}

function navigateToAlbums() {
	state.currentView = 'albums';
	state.currentAlbum = null;
	showView('albums-view');
	history.pushState(null, '', '/');
}

async function navigateToAlbumDetails(albumId) {
	try {
		state.currentAlbum = await GetAlbum(albumId);
		renderAlbumDetails(state.currentAlbum, document.getElementById('album-info'), document.getElementById('tracks-list'));
		state.currentView = 'album-detail';
		showView('album-detail-view');
		history.pushState(null, '', `/album/${albumId}`);
	} catch (error) {
		console.error('Failed to load album details:', error);
		alert('Failed to load album details. Please try again.');
	} 
}

export function showView(viewId) {
	document.querySelectorAll('.view').forEach(view => {
		view.style.display = 'none';
	});
	document.getElementById(viewId).style.display = 'block';
}

document.addEventListener('DOMContentLoaded', init);

window.addEventListener('popstate', () => {
	if (state.currentView === 'album-detail') {
		navigateToAlbums();
	}
});