import { 
	GetCatalog, 
	GetDirPath, 
	NewPlaylist, 
	Play, 

} from '../wailsjs/go/main/App';
import { renderAlbumsList, renderAlbumDetails, updateTrackList } from './album';
import { catalog } from '../wailsjs/go/models';
import { EventsOn, EventsEmit } from '../wailsjs/runtime';
import { formatDuration } from './utils';

interface AppState {
	currentView: 'albums' | 'album-detail';
	currentAlbum: catalog.Album | null;
	currentPlaylistID: string | null;
	currentPlaylistPosistion: number | null;
	catalog: catalog.Catalog | null;
	currentTrack: catalog.Track | null;
	timeElapsed: number | null;
	isPlaying: boolean;
}

export let state: AppState = {
	currentView: 'albums',
	currentAlbum: null,
	currentPlaylistID: null,
	currentPlaylistPosistion: null,
	catalog: null,
	currentTrack: null,
	timeElapsed: null,
	isPlaying: false
};

async function init(): Promise<void> {
	document.getElementById('select-dir-button')?.addEventListener('click', handleSelectDirectory);
	document.getElementById('back-button')?.addEventListener('click', navigateToAlbums);
	document.getElementById('albums-container')?.addEventListener('click', handleAlbumClick);
	document.getElementById('tracks-list')?.addEventListener('click', handleTrackClick);
	document.getElementById('play-button')?.addEventListener('click', handlePlayClick);
	document.getElementById('pause_resume-button')?.addEventListener('click', handlePauseResumeClick)
	document.getElementById('previous_track-button')?.addEventListener('click', handlePreviousTrackClick);
	document.getElementById('next_track-button')?.addEventListener('click', handleNextTrackClick);
	document.getElementById('previous_track-button')?.addEventListener('click', handlePreviousTrackClick);
	document.getElementById('progress-container')?.addEventListener('click', handleProgressBarClick);

	EventsOn("track:playing", (track, playlistPosition) => {
		state.currentTrack = track;
		state.currentPlaylistPosistion = playlistPosition;
		updateNowPlayingUI();
	})

	EventsOn("second:passed", (timeElapsed, playlistID) => {
		if (state.currentPlaylistID === playlistID)
			state.timeElapsed = timeElapsed;
		updateTimeControlUI();
	})

	loadAlbums("");
}

async function handleProgressBarClick(event: MouseEvent): Promise<void> {
	let progBar = document.getElementById('progress-bar') as HTMLDivElement
	let leftX = progBar.getBoundingClientRect().left
	let spanX = progBar.getBoundingClientRect().right - leftX
	let userX = event.clientX
	let ratio = (userX-leftX)/spanX
	let fill = document.getElementById("progress-bar-fill") as HTMLDivElement
	if (!fill)
		return;
	fill.style.width = `${ratio}%`
	if (!state.currentPlaylistID)
		return;
	EventsEmit("ctrl:request", "changePosition", state.currentPlaylistID,`${ratio}`)
}

async function updateTimeControlUI(): Promise<void> {
	let elapsedTimeElement = document.querySelector('.progress-time-elapsed');
	let totalTimeElement = document.querySelector('.progress-time-total');
	let progressBarFill = document.querySelector('.progress-bar-fill');

	if (!elapsedTimeElement || !totalTimeElement || !progressBarFill ||
		!state.currentTrack || state.timeElapsed === null)
		return;

	elapsedTimeElement.textContent = formatDuration(state.timeElapsed);

	if (totalTimeElement.textContent === '0:00') {
		totalTimeElement.textContent = formatDuration(state.currentTrack.Length);
	}

	const progress = (state.timeElapsed / state.currentTrack.Length) * 100;
	(progressBarFill as HTMLElement).style.width = `${progress}%`;
}


async function updateNowPlayingUI(): Promise<void> {
	let trackElement = document.getElementById("current-track");
	let totalTimeElement = document.querySelector('.progress-time-total');
	let elapsedTimeElement = document.querySelector('.progress-time-elapsed');
	let progressBarFill = document.querySelector('.progress-bar-fill');

	if (!trackElement || !state.currentTrack)
		return;

	if (totalTimeElement) totalTimeElement.textContent = '0:00';
	if (elapsedTimeElement) elapsedTimeElement.textContent = '0:00';
	if (progressBarFill) (progressBarFill as HTMLElement).style.width = '0%';

	trackElement.innerHTML = `
        <div class="track-cover">
            <img src="data:image/jpeg;base64,${state.catalog?.Albums[state.currentTrack.AlbumID].Cover || ''}" alt="${state.currentTrack.Title || 'Album cover'}">
        </div>
        <div class="track-details">
            <div class="track-title">${state.currentTrack.Title || 'Unknown Title'}</div>
            <div class="track-artist">${state.currentTrack.Artist?.join(', ') || 'Unknown Artist'}</div>
        </div>
    `;

	if (!state.currentAlbum)
		return;

	if (state.currentView === "album-detail") {
		updateTrackList(state.currentAlbum);
	}
}


async function loadAlbums(fp: string): Promise<void> {
	try {
		state.catalog = await GetCatalog(fp);
		const albumsContainer = document.getElementById('albums-container');
		if (albumsContainer)
			renderAlbumsList(state.catalog.Albums, albumsContainer);
		navigateToAlbums();
	} catch (error) {
		console.error('Failed to load albums:', error);
		alert('Failed to load albums. Please try again.');
	}
}

async function handleSelectDirectory(): Promise<void> {
	try {
		const dirPath = await GetDirPath();
		if (dirPath)
			await loadAlbums(dirPath);
	} catch (error) {
		console.error('Error selecting directory:', error);
		alert(`Failed to select directory: ${error}`);
	}
}

async function handleAlbumClick(event: MouseEvent): Promise<void> {
	const target = event.target as HTMLElement;
	const albumCard = target.closest('.album-card') as HTMLElement;
	if (albumCard && albumCard.dataset.ID)
		navigateToAlbumDetails(albumCard.dataset.ID);
}

async function handlePlayClick(): Promise<void> {
	if (!state.currentAlbum) 
		return;
	let playlist = await NewPlaylist(state.currentAlbum.Tracks);
	state.currentPlaylistID = playlist.ID;
	let nextBtn = document.getElementById('next_track-button') as HTMLButtonElement
	let prevBtn = document.getElementById('previous_track-button') as HTMLButtonElement
	let progBar = document.getElementById('progress-bar') as HTMLDivElement
	nextBtn.className = "playback-ctrl"
	prevBtn.className = "playback-ctrl"
	progBar.className = "progress-bar"
	changePauseResumeButtonState("pause")
	state.isPlaying = true;
	await Play(state.currentPlaylistID);
}

function handlePauseResumeClick(): void {
	if (!state.currentPlaylistID)
		return;
	if (state.isPlaying) {
		changePauseResumeButtonState("resume")
		EventsEmit("ctrl:request", "pause", state.currentPlaylistID)
	} else {
		changePauseResumeButtonState("pause")
		EventsEmit("ctrl:request", "resume", state.currentPlaylistID)
	}
	state.isPlaying = !state.isPlaying
}

function handlePreviousTrackClick(): void {
	if (!state.currentPlaylistID) 
		return;
	changePauseResumeButtonState("pause")
	state.isPlaying = true;
	EventsEmit("ctrl:request", "previous", state.currentPlaylistID)
}

function handleNextTrackClick(): void {
	if (!state.currentPlaylistID) 
		return;
	changePauseResumeButtonState("pause")
	state.isPlaying = true;
	EventsEmit("ctrl:request", "next", state.currentPlaylistID)
}

function changePauseResumeButtonState(to: "pause" | "resume"): void {
	let btn = document.getElementById('pause_resume-button') as HTMLButtonElement
	if (to == "pause") {
		btn.innerHTML = `<i class="material-icons">pause_circle</i>`
		btn.className = "playback-ctrl"
	} else if (to == "resume") {
		btn.innerHTML = `<i class="material-icons">play_circle</i>`
		btn.className = "playback-ctrl"
	}
}	

function handleTrackClick(event: MouseEvent): void {
	const target = event.target as HTMLElement;
	const trackItem = target.closest('.track-item') as HTMLDivElement;
	if (!trackItem || !trackItem.dataset.index) 
		return;
	const trackNumber = trackItem.dataset.index
	if (!state.currentAlbum)
		return
	changePauseResumeButtonState("pause")
	state.isPlaying = true;
	const currentAlbumID = state.currentAlbum.ID
	EventsEmit("track:switch", state.currentPlaylistID, currentAlbumID, trackNumber)
}

function navigateToAlbums(): void {
	state.currentView = 'albums';
	state.currentAlbum = null;
	showView('albums-view');
	history.pushState(null, '', '/');
}

async function navigateToAlbumDetails(albumID: string): Promise<void> {
	try {
		if (state.catalog?.Albums[albumID] == null)
			return;
		state.currentAlbum = state.catalog?.Albums[albumID];
		const albumInfoElement = document.getElementById('album-info');
		const tracksListElement = document.getElementById('tracks-list');

		if (albumInfoElement && tracksListElement) {
			renderAlbumDetails(state.currentAlbum, albumInfoElement, tracksListElement);
		}

		state.currentView = 'album-detail';
		showView('album-detail-view');
		history.pushState(null, '', `/album/${albumID}`);
		document.getElementById('play-button')?.addEventListener('click', handlePlayClick);
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
	if (viewElement)
		viewElement.style.display = 'block';
}

document.addEventListener('DOMContentLoaded', init);

window.addEventListener('popstate', () => {
	if (state.currentView === 'album-detail')
		navigateToAlbums();
});

window.addEventListener('keydown', (e) => {
	let key = e.key.toLowerCase()
	if (key == " " || key == "k") {
		e.preventDefault()
		handlePauseResumeClick()
	} else if (key == "j") {
		e.preventDefault()
		handlePreviousTrackClick()
	} else if (key == "l") {
		e.preventDefault()
		handleNextTrackClick()
	}
})