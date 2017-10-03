import { Injectable } from '@angular/core'
import { Observable } from "rxjs";
import { APIService, Serializable } from "./api";
import 'rxjs/add/operator/map'

export class Feed {
    id: number
    title: string
    description: string
    link: string
    updateError: string
    subscribeError: string
}

export interface OPMLimport {
    opml: string
    dryRun: boolean
}

export class AddFeedResponse extends Serializable {
    success: boolean
    errors: AddFeedError[]
}

export class AddFeedError {
    link: string
    title: string
    error: string
}

class FeedsResponse extends Serializable {
    feeds: Feed[]
}

class OPMLResponse extends Serializable {
    opml: string
}

interface AddFeedData {
    links: string[]
}

@Injectable()
export class FeedService {
    constructor(private api: APIService) { }

    getFeeds() : Observable<Feed[]> {
        return this.api.get("feed").map(response =>
             new FeedsResponse().fromJSON(response.json()).feeds
        );
    }

    discover(query: string) : Observable<Feed[]> {
        return this.api.get(`feed/discover?query=${query}`).map(response =>
             new FeedsResponse().fromJSON(response.json()).feeds
        );
    }

    importOPML(data: OPMLimport): Observable<Feed[]> {
        return this.api.post("opml", JSON.stringify(data)).map(response =>
             new FeedsResponse().fromJSON(response.json()).feeds
        );
    }

    exportOPML(): Observable<string> {
        return this.api.get("opml").map(response =>
             new OPMLResponse().fromJSON(response.json()).opml
        );
    }

    addFeeds(links: string[]) : Observable<AddFeedResponse> {
        let data : AddFeedData = {links: links}
        return this.api.post("feed", JSON.stringify(data)).map(response =>
            new AddFeedResponse().fromJSON(response.json())
        )
    }

    deleteFeed(id: number) : Observable<boolean> {
        return this.api.delete(`feed/${id}`).map(response =>
            !!response.json()["success"]
        )
    }

    updateTags(id: number, tags: string[]) : Observable<boolean> {
        return this.api.put(`feed/${id}/tags`, JSON.stringify(tags)).map(response =>
            !!response.json()["success"]
        )
    }
}