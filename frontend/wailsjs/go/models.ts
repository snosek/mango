export namespace catalog {
	
	export class Track {
	    Title: string;
	    Artist: string[];
	    TrackNumber: number;
	    Length: number;
	    SampleRate: number;
	    Cover?: string;
	    Filepath: string;
	
	    static createFrom(source: any = {}) {
	        return new Track(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Title = source["Title"];
	        this.Artist = source["Artist"];
	        this.TrackNumber = source["TrackNumber"];
	        this.Length = source["Length"];
	        this.SampleRate = source["SampleRate"];
	        this.Cover = source["Cover"];
	        this.Filepath = source["Filepath"];
	    }
	}
	export class Album {
	    Title: string;
	    Artist: string[];
	    Genre: string[];
	    Length: number;
	    Tracks: Track[];
	    Cover: string;
	    Filepath: string;
	
	    static createFrom(source: any = {}) {
	        return new Album(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Title = source["Title"];
	        this.Artist = source["Artist"];
	        this.Genre = source["Genre"];
	        this.Length = source["Length"];
	        this.Tracks = this.convertValues(source["Tracks"], Track);
	        this.Cover = source["Cover"];
	        this.Filepath = source["Filepath"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Catalog {
	    Albums: Album[];
	    Filepath: string;
	
	    static createFrom(source: any = {}) {
	        return new Catalog(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Albums = this.convertValues(source["Albums"], Album);
	        this.Filepath = source["Filepath"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace player {
	
	export class Player {
	
	
	    static createFrom(source: any = {}) {
	        return new Player(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class Playlist {
	    Tracks: catalog.Track[];
	    Current: number;
	    // Go type: Player
	    Player?: any;
	
	    static createFrom(source: any = {}) {
	        return new Playlist(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Tracks = this.convertValues(source["Tracks"], catalog.Track);
	        this.Current = source["Current"];
	        this.Player = this.convertValues(source["Player"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

