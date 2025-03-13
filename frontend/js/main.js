import { GetCatalog, GetAlbum, GetDirPath } from '../wailsjs/go/main/App';
import { catalog } from '../wailsjs/go/models.ts';
import { renderAlbumsList, renderAlbumDetails } from './album.js';
import { playSong } from './player.js';

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
	playSong(state.currentAlbum.Tracks[0].Filepath)
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