import { Component, OnInit, OnDestroy } from "@angular/core";
import { Router } from '@angular/router';
import { Subscription } from "rxjs";
import { listRoute, articleRoute } from "./routing-util"
import { TokenService } from "../services/auth";
import { filter, switchMap, map, combineLatest } from "rxjs/operators";

@Component({
    templateUrl: "./main.html",
    styleUrls: ["./main.css"],
})
export class MainComponent implements OnInit, OnDestroy {
    showsArticle = false
    inSearch = false

    private subscription : Subscription;

    constructor(
        private tokenService: TokenService,
        private router: Router,
    ) {}

    ngOnInit() {
        this.subscription = this.tokenService.tokenObservable().pipe(
            filter(token => token != ""),
            switchMap(token =>
                articleRoute(this.router).pipe(
                    map(route => route != null),
                    combineLatest(
                        listRoute(this.router).pipe(map(route =>
                            route != null && route.data["primary"] == "search"
                        )),
                        (inArticles, inSearch) : [boolean, boolean] =>
                            [inArticles, inSearch]
                    ),
                ),
            ),
        ).subscribe(
            data => {
                this.showsArticle = data[0];
                this.inSearch = data[1];
            },
            error => console.log(error)
        );
    }

    ngOnDestroy(): void {
        this.subscription.unsubscribe()
    }
}
