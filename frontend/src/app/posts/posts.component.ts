import {Component, OnInit, Provider} from '@angular/core';
import {ProviderService} from "../shared/services/provider.service";
import {Tweet} from "../shared/models/models";

@Component({
  selector: 'app-posts',
  templateUrl: './posts.component.html',
  styleUrls: ['./posts.component.css']
})
export class PostsComponent implements OnInit {

  public postContent = '';
  public tweets: Tweet[] = [];

  constructor(private service: ProviderService) { }

  ngOnInit(): void {
    this.service.getFeed().then(res => {
      console.log(res)
      this.tweets = res
    })
  }

  createPost(): void {
    if (this.postContent !== ''){
      console.log(this.postContent)
      this.service.createPost(this.postContent)
    }
  }
}
