import { formatDuration } from './utils';
import { catalog } from '../wailsjs/go/models';

export function renderAlbumsList(albums: Record<string, catalog.Album> | undefined, container: HTMLElement): void {
	container.innerHTML = '';
	if (!albums || Object.keys(albums).length === 0) {
		container.innerHTML = '<p class="no-albums">No albums found. Click "Add Music Directory" to add your music.</p>';
		return;
	}
	const fragment = document.createDocumentFragment();
	for (const albumID in albums) {
		const albumElement = createAlbumCard(albums[albumID]);
		fragment.appendChild(albumElement);
	}
	container.appendChild(fragment);
}

function createAlbumCard(album: catalog.Album): HTMLElement {
	const element = document.createElement('div');
	element.className = 'album-card';
	element.dataset.ID = album.ID;
	element.innerHTML = `
    <div class="album-card__cover">
      <img src="data:image/jpeg;base64,${album.Cover}" alt="${album.Title}"/>
    </div>
    <div class="album-card__info clickable">
      <h3 class="album-card__title clickable">${album.Title || 'Unknown Album'}</h3>
      <p class="album-card__artist clickable">${album.Artist?.join(', ') || 'Unknown Artist'}</p>
    </div>
  `;
	return element;
}

export function renderAlbumDetails(album: catalog.Album, infoContainer: HTMLElement, tracksContainer: HTMLElement): void {
	infoContainer.innerHTML = `
    <h2 class="album-info__title">${album.Title || 'Unknown Album'}</h2>
    <p class="album-info__artist">${album.Artist.join(', ') || 'Unknown Artist'}</p>
    <div class="album-info__details">
      <div>Tracks: ${album.Tracks.length}</div>
    </div>
	<button class="button" id="play-button">Play album</button>
  `;
	tracksContainer.innerHTML = '';
	const fragment = document.createDocumentFragment();
	album.Tracks.forEach((track, index) => {
		fragment.appendChild(createTrackItem(track, index));
	});
	tracksContainer.appendChild(fragment);
}

function createTrackItem(track: catalog.Track, index: number): HTMLElement {
	const element = document.createElement('div');
	element.className = 'track-item';
	element.dataset.index = index.toString();
	element.dataset.filepath = track.Filepath;
	const trackNumber = track.TrackNumber > 0 ? track.TrackNumber : index + 1;
	const duration = formatDuration(track.Length);
	element.innerHTML = `
    <span class="track-item__number">${trackNumber}</span>
    <span class="track-item__title">${track.Title || 'Unknown Title'}</span>
    <span class="track-item__duration">${duration}</span>
  `;
	return element;
}