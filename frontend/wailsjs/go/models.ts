export namespace catalog {
	
	export class Track {
	    Filepath: string;
	    Title: string;
	    Artist: string[];
	    Genre: string[];
	    TrackNumber: number;
	    Length: number;
	    SampleRate: number;
	
	    static createFrom(source: any = {}) {
	        return new Track(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Filepath = source["Filepath"];
	        this.Title = source["Title"];
	        this.Artist = source["Artist"];
	        this.Genre = source["Genre"];
	        this.TrackNumber = source["TrackNumber"];
	        this.Length = source["Length"];
	        this.SampleRate = source["SampleRate"];
	    }
	}
	export class AlbumMetadata {
	    Filepath: string;
	    Title: string;
	    Artist: string[];
	    Genre: string[];
	    Length: number;
	    SampleRate: number;
	
	    static createFrom(source: any = {}) {
	        return new AlbumMetadata(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Filepath = source["Filepath"];
	        this.Title = source["Title"];
	        this.Artist = source["Artist"];
	        this.Genre = source["Genre"];
	        this.Length = source["Length"];
	        this.SampleRate = source["SampleRate"];
	    }
	}
	export class Album {
	    Metadata?: AlbumMetadata;
	    Tracks: Track[];
	
	    static createFrom(source: any = {}) {
	        return new Album(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Metadata = this.convertValues(source["Metadata"], AlbumMetadata);
	        this.Tracks = this.convertValues(source["Tracks"], Track);
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

