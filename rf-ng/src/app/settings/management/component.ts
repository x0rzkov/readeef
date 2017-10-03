import { Component, OnInit, Inject } from "@angular/core" ;
import { MdDialog, MD_DIALOG_DATA } from '@angular/material';
import { FeedService, Feed } from "../../services/feed";
import { TagService } from "../../services/tag";
import { FaviconService } from "../../services/favicon";
import { Observable } from "rxjs";
import 'rxjs/add/operator/combineLatest'

@Component({
    selector: "settings-management",
    templateUrl: "./management.html",
    styleUrls: ["../common.css", "./management.css"]
})
export class ManagementSettingsComponent implements OnInit {
    feeds = new Array<[Feed, string[]]>()

    constructor(
        private feedService: FeedService,
        private tagService: TagService,
        private faviconService: FaviconService,
        private errorDialog: MdDialog,
    ) {}

    ngOnInit(): void {
        this.feedService.getFeeds().combineLatest(
            this.tagService.getTagsFeedIDs(),
            (feeds, tags) : [Feed, string[]][] => {
                let tagMap  = new Map<number, string[]>()
                for (let tag of tags) {
                    for (let id of tag.ids) {
                        if (tagMap.has(id)) {
                            tagMap.get(id).push(tag.tag.value);
                        } else {
                            tagMap.set(id, [tag.tag.value]);
                        }
                    }
                }

                return feeds.map((feed) : [Feed, string[]] => 
                    [feed, tagMap.get(feed.id) || []]
                );
            }
        ).subscribe(
            feeds => this.feeds = feeds || [],
            error => console.log(error),
        );
    }

    favicon(url: string) : string {
        return this.faviconService.iconURL(url);
    }

    tagsChange(event: Event, feedID: number) {
        let tags : string = event.target["value"];

        this.feedService.updateTags(
            feedID, tags.split(',').map(tag => tag.trim())
        ).subscribe(
            success => {},
            error => console.log(error),
        );
    }

    showError(error: string) {
        this.errorDialog.open(ErrorDialog, {
            width: "300px",
             data: error.split("\n").filter(err => err),
        });
    }

    deleteFeed(event: Event, feedID: number) {
        this.feedService.deleteFeed(
            feedID
        ).subscribe(
            success => {
                if (success) {
                    let el = event.target["parentNode"];
                    while ((el = el.parentElement) && !el.classList.contains("feed"));
                    el.parentNode.removeChild(el);
                }
            },
            error => console.log(error),
        );
    }
}

@Component({
    selector: 'error-dialog',
    templateUrl: 'error-dialog.html',
})
export class ErrorDialog {
    constructor(@Inject(MD_DIALOG_DATA) public errors: string[]) {}
}