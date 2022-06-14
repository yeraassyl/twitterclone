import { Component, OnInit } from '@angular/core';
import {User} from "../shared/models/models";
import {ProviderService} from "../shared/services/provider.service";

@Component({
  selector: 'app-users',
  templateUrl: './users.component.html',
  styleUrls: ['./users.component.css']
})
export class UsersComponent implements OnInit {

  public users: User[] = [];

  constructor(private provider: ProviderService) { }

  ngOnInit(): void {
    this.provider.listUsers().then(res => this.users = res);
  }

  follow(id: number): void {
    this.provider.follow(id);
  }
}
