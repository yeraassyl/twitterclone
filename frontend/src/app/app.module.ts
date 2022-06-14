import {ClassProvider, NgModule} from '@angular/core';
import {BrowserModule} from '@angular/platform-browser';

import {AppComponent} from './app.component';
import {MainComponent} from './main/main.component';
import {HTTP_INTERCEPTORS, HttpClientModule} from "@angular/common/http";
import {ProviderService} from "./shared/services/provider.service";

import {
  SocialLoginModule,
  SocialAuthServiceConfig,
  GoogleLoginProvider
} from 'angularx-social-login';
import {AuthInterceptor} from "./AuthInterceptor";
import { PostsComponent } from './posts/posts.component';
import { UserComponent } from './user/user.component';
import { UsersComponent } from './users/users.component';
import { AppRoutingModule } from './app-routing.module';
import {FormsModule} from "@angular/forms";
import { PostComponent } from './post/post.component';

@NgModule({
  declarations: [
    AppComponent,
    MainComponent,
    PostsComponent,
    UserComponent,
    UsersComponent,
    PostComponent
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    SocialLoginModule,
    AppRoutingModule,
    FormsModule,
  ],
  providers: [
    ProviderService,
    {
      provide: 'SocialAuthServiceConfig',
      useValue: {
        autoLogin: false,
        providers: [
          {
            id: GoogleLoginProvider.PROVIDER_ID,
            provider: new GoogleLoginProvider(
              '798066806591-sn722ltj9mus74s6985moee0mq9cnl0u.apps.googleusercontent.com'
            )
          }
        ]
      } as SocialAuthServiceConfig
    },
    {
      provide: HTTP_INTERCEPTORS,
      useClass: AuthInterceptor,
      multi: true
    } as ClassProvider
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}
