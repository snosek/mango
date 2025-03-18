export function formatDuration(nanoseconds: number): string {
	const seconds = Math.floor(nanoseconds / 1000000000);
	const minutes = Math.floor(seconds / 60);
	const remainingSeconds = seconds % 60;
	return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
}