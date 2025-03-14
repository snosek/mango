import { PlaySong, PauseSong } from "../wailsjs/go/main/App";

export function playSong(fp) {
	PlaySong(fp)
}

export function pauseSong() {
	PauseSong()
}