import {Component, OnInit} from '@angular/core';
import {ProviderService} from "../shared/services/provider.service";
import {GoogleLoginProvider, SocialAuthService, SocialUser} from "angularx-social-login";
import {Tweet, User, UserSimple, UserSimplImpl} from "../shared/models/models";

@Component({
  selector: 'app-main',
  templateUrl: './main.component.html',
  styleUrls: ['./main.component.css']
})
export class MainComponent implements OnInit {
  get myTweets(): Tweet[] {
    return this._myTweets;
  }
  set myTweets(value: Tweet[]) {
    this._myTweets = value;
  }
  public email: string | null = ''
  public logged = false
  public socialUser!: SocialUser;
  public followers: UserSimple[] = [];
  public following: UserSimple[] = [];
  private _myTweets: Tweet[];
  public user!: User;

  constructor(private provider: ProviderService,
              private socialAuthService: SocialAuthService) {
    this._myTweets = [];
  }

  ngOnInit(): void {
    if (localStorage.getItem('token') != null) {
      this.logged = true
      this.email = localStorage.getItem('email')
      this.provider.createUserIfNotExists().then(res => {
        this.user = res;
        console.log(res.following);
      });
      this.getFollowers();
      this.getFollowing();
      this.myTweets = this.getMyTweets();
    }
  }

  login() {
    this.socialAuthService.signIn(GoogleLoginProvider.PROVIDER_ID).then(user => {
      this.socialUser = user;
      localStorage.setItem('token', user.idToken)
      localStorage.setItem('email', user.email)
      this.email = user.email
      this.provider.createUserIfNotExists().then(res => this.user = res);
    })
  }

  logOut() {
    this.socialAuthService.signOut();
    localStorage.removeItem('token')
    localStorage.removeItem('email')
    this.logged = false;
  }

  userList() {
    this.provider.listUsers();
  }

  getFollowers() {
    this.user.followers.forEach(id => {
      this.provider.getUser(id).then(user => {
        console.log(user.id);
        this.followers.push(new UserSimplImpl(user.email, user.id))
      })
    })
  }

  getFollowing() {
    this.user.following.forEach(id => {
      this.provider.getUser(id).then(user => {
        console.log(user.id);
        this.following.push(new UserSimplImpl(user.email, user.id))
      })
    })
  }

  getMyTweets(): Tweet[] {
    let tweets: Tweet[] = [];
    this.provider.getPosts(this.user.id).then(res => tweets = res);
    return tweets;
  }
}
