import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { HttpModule, BaseRequestOptions } from '@angular/http';

import { AppComponent } from './components/app';
import { AppRouting } from './app.routing';

import { MainModule } from './main/module';
import { LoginModule } from './login/module';

import { AuthGuard } from './guards/auth';

import { TokenService } from './services/auth';
import { APIService } from './services/api';
import { EventService } from './services/events';
import { FeaturesService } from './services/features';
import { ArticleService } from './services/article';
import { FeedService } from './services/feed';
import { TagService } from './services/tag';
import { PreferencesService } from './services/preferences';

@NgModule({
  declarations: [
    AppComponent,
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    HttpModule,
    AppRouting,
    LoginModule,
    MainModule,
  ],
  providers: [
    TokenService,
    APIService,
    EventService,
    FeaturesService,
    ArticleService,
    FeedService,
    TagService,
    PreferencesService,
    AuthGuard,
    BaseRequestOptions
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
