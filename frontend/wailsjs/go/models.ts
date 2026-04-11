export namespace backend {
	
	export class AnalysisDecodeResponse {
	    pcm_base64: string;
	    sample_rate: number;
	    channels: number;
	    bits_per_sample: number;
	    duration: number;
	    bitrate_kbps?: number;
	    bit_depth?: string;
	
	    static createFrom(source: any = {}) {
	        return new AnalysisDecodeResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pcm_base64 = source["pcm_base64"];
	        this.sample_rate = source["sample_rate"];
	        this.channels = source["channels"];
	        this.bits_per_sample = source["bits_per_sample"];
	        this.duration = source["duration"];
	        this.bitrate_kbps = source["bitrate_kbps"];
	        this.bit_depth = source["bit_depth"];
	    }
	}
	export class AudioMetadata {
	    title: string;
	    artist: string;
	    album: string;
	    album_artist: string;
	    track_number: number;
	    disc_number: number;
	    year: string;
	
	    static createFrom(source: any = {}) {
	        return new AudioMetadata(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.artist = source["artist"];
	        this.album = source["album"];
	        this.album_artist = source["album_artist"];
	        this.track_number = source["track_number"];
	        this.disc_number = source["disc_number"];
	        this.year = source["year"];
	    }
	}
	export class AvatarDownloadResponse {
	    success: boolean;
	    message: string;
	    file?: string;
	    error?: string;
	    already_exists?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AvatarDownloadResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.file = source["file"];
	        this.error = source["error"];
	        this.already_exists = source["already_exists"];
	    }
	}
	export class ConvertAudioResult {
	    input_file: string;
	    output_file: string;
	    success: boolean;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new ConvertAudioResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.input_file = source["input_file"];
	        this.output_file = source["output_file"];
	        this.success = source["success"];
	        this.error = source["error"];
	    }
	}
	export class CoverDownloadResponse {
	    success: boolean;
	    message: string;
	    file?: string;
	    error?: string;
	    already_exists?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CoverDownloadResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.file = source["file"];
	        this.error = source["error"];
	        this.already_exists = source["already_exists"];
	    }
	}
	export class DownloadItem {
	    id: string;
	    track_name: string;
	    artist_name: string;
	    album_name: string;
	    spotify_id: string;
	    status: string;
	    progress: number;
	    total_size: number;
	    speed: number;
	    start_time: number;
	    end_time: number;
	    error_message: string;
	    file_path: string;
	
	    static createFrom(source: any = {}) {
	        return new DownloadItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.track_name = source["track_name"];
	        this.artist_name = source["artist_name"];
	        this.album_name = source["album_name"];
	        this.spotify_id = source["spotify_id"];
	        this.status = source["status"];
	        this.progress = source["progress"];
	        this.total_size = source["total_size"];
	        this.speed = source["speed"];
	        this.start_time = source["start_time"];
	        this.end_time = source["end_time"];
	        this.error_message = source["error_message"];
	        this.file_path = source["file_path"];
	    }
	}
	export class DownloadQueueInfo {
	    is_downloading: boolean;
	    queue: DownloadItem[];
	    current_speed: number;
	    total_downloaded: number;
	    session_start_time: number;
	    queued_count: number;
	    completed_count: number;
	    failed_count: number;
	    skipped_count: number;
	
	    static createFrom(source: any = {}) {
	        return new DownloadQueueInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.is_downloading = source["is_downloading"];
	        this.queue = this.convertValues(source["queue"], DownloadItem);
	        this.current_speed = source["current_speed"];
	        this.total_downloaded = source["total_downloaded"];
	        this.session_start_time = source["session_start_time"];
	        this.queued_count = source["queued_count"];
	        this.completed_count = source["completed_count"];
	        this.failed_count = source["failed_count"];
	        this.skipped_count = source["skipped_count"];
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
	export class FetchHistoryItem {
	    id: string;
	    url: string;
	    type: string;
	    name: string;
	    info: string;
	    image: string;
	    data: string;
	    timestamp: number;
	
	    static createFrom(source: any = {}) {
	        return new FetchHistoryItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.url = source["url"];
	        this.type = source["type"];
	        this.name = source["name"];
	        this.info = source["info"];
	        this.image = source["image"];
	        this.data = source["data"];
	        this.timestamp = source["timestamp"];
	    }
	}
	export class FileInfo {
	    name: string;
	    path: string;
	    is_dir: boolean;
	    size: number;
	    children?: FileInfo[];
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.is_dir = source["is_dir"];
	        this.size = source["size"];
	        this.children = this.convertValues(source["children"], FileInfo);
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
	export class FlacInfo {
	    path: string;
	    sample_rate: number;
	    bits_per_sample: number;
	
	    static createFrom(source: any = {}) {
	        return new FlacInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.sample_rate = source["sample_rate"];
	        this.bits_per_sample = source["bits_per_sample"];
	    }
	}
	export class GalleryImageDownloadResponse {
	    success: boolean;
	    message: string;
	    file?: string;
	    error?: string;
	    already_exists?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new GalleryImageDownloadResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.file = source["file"];
	        this.error = source["error"];
	        this.already_exists = source["already_exists"];
	    }
	}
	export class HeaderDownloadResponse {
	    success: boolean;
	    message: string;
	    file?: string;
	    error?: string;
	    already_exists?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new HeaderDownloadResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.file = source["file"];
	        this.error = source["error"];
	        this.already_exists = source["already_exists"];
	    }
	}
	export class HistoryItem {
	    id: string;
	    spotify_id: string;
	    title: string;
	    artists: string;
	    album: string;
	    duration_str: string;
	    cover_url: string;
	    quality: string;
	    format: string;
	    path: string;
	    source: string;
	    timestamp: number;
	
	    static createFrom(source: any = {}) {
	        return new HistoryItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.spotify_id = source["spotify_id"];
	        this.title = source["title"];
	        this.artists = source["artists"];
	        this.album = source["album"];
	        this.duration_str = source["duration_str"];
	        this.cover_url = source["cover_url"];
	        this.quality = source["quality"];
	        this.format = source["format"];
	        this.path = source["path"];
	        this.source = source["source"];
	        this.timestamp = source["timestamp"];
	    }
	}
	export class LyricsDownloadResponse {
	    success: boolean;
	    message: string;
	    file?: string;
	    error?: string;
	    already_exists?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new LyricsDownloadResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.file = source["file"];
	        this.error = source["error"];
	        this.already_exists = source["already_exists"];
	    }
	}
	export class ProgressInfo {
	    is_downloading: boolean;
	    mb_downloaded: number;
	    speed_mbps: number;
	
	    static createFrom(source: any = {}) {
	        return new ProgressInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.is_downloading = source["is_downloading"];
	        this.mb_downloaded = source["mb_downloaded"];
	        this.speed_mbps = source["speed_mbps"];
	    }
	}
	export class RenamePreview {
	    old_path: string;
	    old_name: string;
	    new_name: string;
	    new_path: string;
	    error?: string;
	    metadata: AudioMetadata;
	
	    static createFrom(source: any = {}) {
	        return new RenamePreview(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.old_path = source["old_path"];
	        this.old_name = source["old_name"];
	        this.new_name = source["new_name"];
	        this.new_path = source["new_path"];
	        this.error = source["error"];
	        this.metadata = this.convertValues(source["metadata"], AudioMetadata);
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
	export class RenameResult {
	    old_path: string;
	    new_path: string;
	    success: boolean;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new RenameResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.old_path = source["old_path"];
	        this.new_path = source["new_path"];
	        this.success = source["success"];
	        this.error = source["error"];
	    }
	}
	export class ResampleResult {
	    input_file: string;
	    output_file: string;
	    success: boolean;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new ResampleResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.input_file = source["input_file"];
	        this.output_file = source["output_file"];
	        this.success = source["success"];
	        this.error = source["error"];
	    }
	}
	export class SearchResult {
	    id: string;
	    name: string;
	    type: string;
	    artists?: string;
	    album_name?: string;
	    images: string;
	    release_date?: string;
	    external_urls: string;
	    duration_ms?: number;
	    total_tracks?: number;
	    owner?: string;
	    is_explicit?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SearchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.artists = source["artists"];
	        this.album_name = source["album_name"];
	        this.images = source["images"];
	        this.release_date = source["release_date"];
	        this.external_urls = source["external_urls"];
	        this.duration_ms = source["duration_ms"];
	        this.total_tracks = source["total_tracks"];
	        this.owner = source["owner"];
	        this.is_explicit = source["is_explicit"];
	    }
	}
	export class SearchResponse {
	    tracks: SearchResult[];
	    albums: SearchResult[];
	    artists: SearchResult[];
	    playlists: SearchResult[];
	
	    static createFrom(source: any = {}) {
	        return new SearchResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tracks = this.convertValues(source["tracks"], SearchResult);
	        this.albums = this.convertValues(source["albums"], SearchResult);
	        this.artists = this.convertValues(source["artists"], SearchResult);
	        this.playlists = this.convertValues(source["playlists"], SearchResult);
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

export namespace main {
	
	export class AvatarDownloadRequest {
	    avatar_url: string;
	    artist_name: string;
	    output_dir: string;
	
	    static createFrom(source: any = {}) {
	        return new AvatarDownloadRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.avatar_url = source["avatar_url"];
	        this.artist_name = source["artist_name"];
	        this.output_dir = source["output_dir"];
	    }
	}
	export class CheckFileExistenceRequest {
	    spotify_id: string;
	    track_name: string;
	    artist_name: string;
	    album_name?: string;
	    album_artist?: string;
	    release_date?: string;
	    track_number?: number;
	    disc_number?: number;
	    position?: number;
	    use_album_track_number?: boolean;
	    filename_format?: string;
	    include_track_number?: boolean;
	    audio_format?: string;
	    relative_path?: string;
	
	    static createFrom(source: any = {}) {
	        return new CheckFileExistenceRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.spotify_id = source["spotify_id"];
	        this.track_name = source["track_name"];
	        this.artist_name = source["artist_name"];
	        this.album_name = source["album_name"];
	        this.album_artist = source["album_artist"];
	        this.release_date = source["release_date"];
	        this.track_number = source["track_number"];
	        this.disc_number = source["disc_number"];
	        this.position = source["position"];
	        this.use_album_track_number = source["use_album_track_number"];
	        this.filename_format = source["filename_format"];
	        this.include_track_number = source["include_track_number"];
	        this.audio_format = source["audio_format"];
	        this.relative_path = source["relative_path"];
	    }
	}
	export class CheckFileExistenceResult {
	    spotify_id: string;
	    exists: boolean;
	    file_path?: string;
	    track_name?: string;
	    artist_name?: string;
	
	    static createFrom(source: any = {}) {
	        return new CheckFileExistenceResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.spotify_id = source["spotify_id"];
	        this.exists = source["exists"];
	        this.file_path = source["file_path"];
	        this.track_name = source["track_name"];
	        this.artist_name = source["artist_name"];
	    }
	}
	export class ConvertAudioRequest {
	    input_files: string[];
	    output_format: string;
	    bitrate: string;
	    codec: string;
	
	    static createFrom(source: any = {}) {
	        return new ConvertAudioRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.input_files = source["input_files"];
	        this.output_format = source["output_format"];
	        this.bitrate = source["bitrate"];
	        this.codec = source["codec"];
	    }
	}
	export class CoverDownloadRequest {
	    cover_url: string;
	    track_name: string;
	    artist_name: string;
	    album_name: string;
	    album_artist: string;
	    release_date: string;
	    output_dir: string;
	    filename_format: string;
	    track_number: boolean;
	    position: number;
	    disc_number: number;
	
	    static createFrom(source: any = {}) {
	        return new CoverDownloadRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cover_url = source["cover_url"];
	        this.track_name = source["track_name"];
	        this.artist_name = source["artist_name"];
	        this.album_name = source["album_name"];
	        this.album_artist = source["album_artist"];
	        this.release_date = source["release_date"];
	        this.output_dir = source["output_dir"];
	        this.filename_format = source["filename_format"];
	        this.track_number = source["track_number"];
	        this.position = source["position"];
	        this.disc_number = source["disc_number"];
	    }
	}
	export class DownloadFFmpegResponse {
	    success: boolean;
	    message: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new DownloadFFmpegResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.error = source["error"];
	    }
	}
	export class DownloadRequest {
	    service: string;
	    query?: string;
	    track_name?: string;
	    artist_name?: string;
	    album_name?: string;
	    album_artist?: string;
	    release_date?: string;
	    cover_url?: string;
	    api_url?: string;
	    output_dir?: string;
	    audio_format?: string;
	    filename_format?: string;
	    track_number?: boolean;
	    position?: number;
	    use_album_track_number?: boolean;
	    spotify_id?: string;
	    embed_lyrics?: boolean;
	    embed_max_quality_cover?: boolean;
	    service_url?: string;
	    duration?: number;
	    item_id?: string;
	    spotify_track_number?: number;
	    spotify_disc_number?: number;
	    spotify_total_tracks?: number;
	    spotify_total_discs?: number;
	    copyright?: string;
	    publisher?: string;
	    playlist_name?: string;
	    playlist_owner?: string;
	    allow_fallback: boolean;
	    use_first_artist_only?: boolean;
	    use_single_genre?: boolean;
	    embed_genre?: boolean;
	    separator?: string;
	
	    static createFrom(source: any = {}) {
	        return new DownloadRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.service = source["service"];
	        this.query = source["query"];
	        this.track_name = source["track_name"];
	        this.artist_name = source["artist_name"];
	        this.album_name = source["album_name"];
	        this.album_artist = source["album_artist"];
	        this.release_date = source["release_date"];
	        this.cover_url = source["cover_url"];
	        this.api_url = source["api_url"];
	        this.output_dir = source["output_dir"];
	        this.audio_format = source["audio_format"];
	        this.filename_format = source["filename_format"];
	        this.track_number = source["track_number"];
	        this.position = source["position"];
	        this.use_album_track_number = source["use_album_track_number"];
	        this.spotify_id = source["spotify_id"];
	        this.embed_lyrics = source["embed_lyrics"];
	        this.embed_max_quality_cover = source["embed_max_quality_cover"];
	        this.service_url = source["service_url"];
	        this.duration = source["duration"];
	        this.item_id = source["item_id"];
	        this.spotify_track_number = source["spotify_track_number"];
	        this.spotify_disc_number = source["spotify_disc_number"];
	        this.spotify_total_tracks = source["spotify_total_tracks"];
	        this.spotify_total_discs = source["spotify_total_discs"];
	        this.copyright = source["copyright"];
	        this.publisher = source["publisher"];
	        this.playlist_name = source["playlist_name"];
	        this.playlist_owner = source["playlist_owner"];
	        this.allow_fallback = source["allow_fallback"];
	        this.use_first_artist_only = source["use_first_artist_only"];
	        this.use_single_genre = source["use_single_genre"];
	        this.embed_genre = source["embed_genre"];
	        this.separator = source["separator"];
	    }
	}
	export class DownloadResponse {
	    success: boolean;
	    message: string;
	    file?: string;
	    error?: string;
	    already_exists?: boolean;
	    item_id?: string;
	
	    static createFrom(source: any = {}) {
	        return new DownloadResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.file = source["file"];
	        this.error = source["error"];
	        this.already_exists = source["already_exists"];
	        this.item_id = source["item_id"];
	    }
	}
	export class GalleryImageDownloadRequest {
	    image_url: string;
	    artist_name: string;
	    image_index: number;
	    output_dir: string;
	
	    static createFrom(source: any = {}) {
	        return new GalleryImageDownloadRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.image_url = source["image_url"];
	        this.artist_name = source["artist_name"];
	        this.image_index = source["image_index"];
	        this.output_dir = source["output_dir"];
	    }
	}
	export class HeaderDownloadRequest {
	    header_url: string;
	    artist_name: string;
	    output_dir: string;
	
	    static createFrom(source: any = {}) {
	        return new HeaderDownloadRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.header_url = source["header_url"];
	        this.artist_name = source["artist_name"];
	        this.output_dir = source["output_dir"];
	    }
	}
	export class InstallFFmpegWithBrewResponse {
	    success: boolean;
	    message: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new InstallFFmpegWithBrewResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.error = source["error"];
	    }
	}
	export class LyricsDownloadRequest {
	    spotify_id: string;
	    track_name: string;
	    artist_name: string;
	    album_name: string;
	    album_artist: string;
	    release_date: string;
	    output_dir: string;
	    filename_format: string;
	    track_number: boolean;
	    position: number;
	    use_album_track_number: boolean;
	    disc_number: number;
	
	    static createFrom(source: any = {}) {
	        return new LyricsDownloadRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.spotify_id = source["spotify_id"];
	        this.track_name = source["track_name"];
	        this.artist_name = source["artist_name"];
	        this.album_name = source["album_name"];
	        this.album_artist = source["album_artist"];
	        this.release_date = source["release_date"];
	        this.output_dir = source["output_dir"];
	        this.filename_format = source["filename_format"];
	        this.track_number = source["track_number"];
	        this.position = source["position"];
	        this.use_album_track_number = source["use_album_track_number"];
	        this.disc_number = source["disc_number"];
	    }
	}
	export class ResampleAudioRequest {
	    input_files: string[];
	    sample_rate: string;
	    bit_depth: string;
	
	    static createFrom(source: any = {}) {
	        return new ResampleAudioRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.input_files = source["input_files"];
	        this.sample_rate = source["sample_rate"];
	        this.bit_depth = source["bit_depth"];
	    }
	}
	export class SpotifyMetadataRequest {
	    url: string;
	    batch: boolean;
	    delay: number;
	    timeout: number;
	    separator?: string;
	
	    static createFrom(source: any = {}) {
	        return new SpotifyMetadataRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.batch = source["batch"];
	        this.delay = source["delay"];
	        this.timeout = source["timeout"];
	        this.separator = source["separator"];
	    }
	}
	export class SpotifySearchByTypeRequest {
	    query: string;
	    search_type: string;
	    limit: number;
	    offset: number;
	
	    static createFrom(source: any = {}) {
	        return new SpotifySearchByTypeRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.query = source["query"];
	        this.search_type = source["search_type"];
	        this.limit = source["limit"];
	        this.offset = source["offset"];
	    }
	}
	export class SpotifySearchRequest {
	    query: string;
	    limit: number;
	
	    static createFrom(source: any = {}) {
	        return new SpotifySearchRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.query = source["query"];
	        this.limit = source["limit"];
	    }
	}

}

