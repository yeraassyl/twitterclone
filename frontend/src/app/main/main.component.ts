import {Component, OnInit} from '@angular/core';
import {ProviderService} from "../shared/services/provider.service";
import {GoogleLoginProvider, SocialAuthService, SocialUser} from "angularx-social-login";

@Component({
  selector: 'app-main',
  templateUrl: './main.component.html',
  styleUrls: ['./main.component.css']
})
export class MainComponent implements OnInit {

  public email = ''
  public logged = false
  public socialUser:SocialUser = new SocialUser();

  constructor(private provider: ProviderService,
              private socialAuthService: SocialAuthService) {
  }

  ngOnInit(): void {
    let token = localStorage.getItem('token');
    if (token != null) {
      this.logged = true;
    }
    this.socialAuthService.authState.subscribe((user) => {
      this.socialUser = user;
      localStorage.setItem('token', user.idToken)
    })
  }


  login(){
    this.socialAuthService.signIn(GoogleLoginProvider.PROVIDER_ID)
    this.logged = true;
  }

  logOut(){
    this.socialAuthService.signOut();
    this.logged = false;
    localStorage.clear();
  }

  userList(){
    this.provider.listUsers();
  }
}
