import { Component, OnInit } from '@angular/core';
import {Tweet, User} from "../shared/models/models";
import {ActivatedRoute} from "@angular/router";
import {ProviderService} from "../shared/services/provider.service";

@Component({
  selector: 'app-user',
  templateUrl: './user.component.html',
  styleUrls: ['./user.component.css']
})
export class UserComponent implements OnInit {

  public user!: User;
  public tweets: Tweet[] = [];
  public pathParams: any;
  constructor(private route: ActivatedRoute, private provider: ProviderService) {
    route.params.subscribe(params => this.pathParams = params);
  }

  ngOnInit(): void {
    this.provider.getUser(this.pathParams.id).then(res => this.user = res);
    this.provider.getPosts(this.pathParams.id).then(res => {
      console.log(res);
      this.tweets = res;
    });
  }

  follow(): void {
    this.provider.follow(this.user.id);
  }

}
