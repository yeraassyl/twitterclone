import { Component, OnInit } from '@angular/core';
import {ProviderService} from "../shared/services/provider.service";
import {ActivatedRoute} from "@angular/router";
import {Tweet} from "../shared/models/models";

@Component({
  selector: 'app-post',
  templateUrl: './post.component.html',
  styleUrls: ['./post.component.css']
})
export class PostComponent implements OnInit {

  public likeCount: number = 0;
  public pathParams: any;
  public tweet!: Tweet;
  constructor(private route: ActivatedRoute, private provider: ProviderService) {
    route.params.subscribe(params => this.pathParams = params);
  }

  ngOnInit(): void {
    this.provider.getPost(this.pathParams.id).then(res => this.tweet = res);
    this.likeCount = this.tweet.likes;
  }

  like(): void {
    this.provider.likePost(this.pathParams.id);
    this.likeCount++;
  }
}
