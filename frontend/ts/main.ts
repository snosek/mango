import { GetCatalog, GetAlbum, GetDirPath, NewPlaylist, Play, PauseSong, ResumeSong } from '../wailsjs/go/main/App';
import { catalog } from '../wailsjs/go/models';
import { renderAlbumsList, renderAlbumDetails } from './album';

interface AppState {
	currentView: 'albums' | 'album-detail';
	currentAlbum: catalog.Album | null;
	currentPlaylistID: string | null;
	catalog: catalog.Catalog | null;
}

let state: AppState = {
	currentView: 'albums',
	currentAlbum: null,
	currentPlaylistID: null,
	catalog: null
};

async function init(): Promise<void> {
	document.getElementById('select-dir-button')?.addEventListener('click', handleSelectDirectory);
	document.getElementById('back-button')?.addEventListener('click', navigateToAlbums);
	document.getElementById('albums-container')?.addEventListener('click', handleAlbumClick);
	document.getElementById('tracks-list')?.addEventListener('click', handleTrackClick);
	document.getElementById('play-button')?.addEventListener('click', handlePlayClick);
	document.getElementById('pause-button')?.addEventListener('click', handlePauseClick);
	document.getElementById('resume-button')?.addEventListener('click', handleResumeClick);

	loadAlbums("");
}

async function loadAlbums(fp: string): Promise<void> {
	try {
		state.catalog = await GetCatalog(fp);
		const albumsContainer = document.getElementById('albums-container');
		if (albumsContainer) {
			renderAlbumsList(state.catalog.Albums, albumsContainer);
		}
		navigateToAlbums();
	} catch (error) {
		console.error('Failed to load albums:', error);
		alert('Failed to load albums. Please try again.');
	}
}

async function handleSelectDirectory(): Promise<void> {
	try {
		const dirPath = await GetDirPath();
		if (dirPath) {
			await loadAlbums(dirPath);
		}
	} catch (error) {
		console.error('Error selecting directory:', error);
		alert(`Failed to select directory: ${error}`);
	}
}

async function handleAlbumClick(event: MouseEvent): Promise<void> {
	const target = event.target as HTMLElement;
	const albumCard = target.closest('.album-card') as HTMLElement;
	if (albumCard && albumCard.dataset.id) {
		navigateToAlbumDetails(albumCard.dataset.id);
	}
}

async function handlePlayClick(): Promise<void> {
	if (!state.currentAlbum) return;

	let playlist = await NewPlaylist(state.currentAlbum.Tracks);
	state.currentPlaylistID = playlist.ID;
	await Play(state.currentPlaylistID);
}

function handlePauseClick(): void {
	if (!state.currentPlaylistID) return;
	PauseSong(state.currentPlaylistID);
}

function handleResumeClick(): void {
	if (!state.currentPlaylistID) return;
	ResumeSong(state.currentPlaylistID);
}

function handleTrackClick(event: MouseEvent): void {
	const target = event.target as HTMLElement;
	const trackItem = target.closest('.track-item');
	if (!trackItem) return;
	// Track click functionality can be added here
}

function navigateToAlbums(): void {
	state.currentView = 'albums';
	state.currentAlbum = null;
	showView('albums-view');
	history.pushState(null, '', '/');
}

async function navigateToAlbumDetails(albumId: string): Promise<void> {
	try {
		state.currentAlbum = await GetAlbum(albumId);
		const albumInfoElement = document.getElementById('album-info');
		const tracksListElement = document.getElementById('tracks-list');

		if (albumInfoElement && tracksListElement) {
			renderAlbumDetails(state.currentAlbum, albumInfoElement, tracksListElement);
		}

		state.currentView = 'album-detail';
		showView('album-detail-view');
		history.pushState(null, '', `/album/${albumId}`);
	} catch (error) {
		console.error('Failed to load album details:', error);
		alert('Failed to load album details. Please try again.');
	}
}

export function showView(viewId: string): void {
	window.scrollTo(0, 0);
	document.querySelectorAll('.view').forEach(view => {
		(view as HTMLElement).style.display = 'none';
	});
	const viewElement = document.getElementById(viewId);
	if (viewElement) {
		viewElement.style.display = 'block';
	}
}

document.addEventListener('DOMContentLoaded', init);

window.addEventListener('popstate', () => {
	if (state.currentView === 'album-detail') {
		navigateToAlbums();
	}
});