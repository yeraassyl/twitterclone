import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {IAuthResponse, Tweet, User} from "../models/models";
import * as moment from "moment";
import {MainService} from "./main.service";

@Injectable({
  providedIn: 'root'
})
export class ProviderService extends MainService{

  private baseUrl = 'http://18.179.197.12'

  constructor(http: HttpClient) {
    super(http)
  }

  formatDate(date: Date) {
    return moment(date).format('YYYY-MM-DD');
  }

  auth(){
    return this.http.get<IAuthResponse>(`http://localhost:8080`).subscribe(resp => {
      console.log(resp);
    })
  }

  listUsers(): Promise<User[]> {
    return this.get(`${this.baseUrl}/user`, {});
  }

  getUser(id: number): Promise<User> {
    return this.get(`${this.baseUrl}/user/${id}`, {})
  }

  createPost(content: string): Promise<any>{
    return this.post(`${this.baseUrl}/tweet/`, {
      content: content
    })
  }

  getFeed(): Promise<Tweet[]> {
    return this.get(`${this.baseUrl}/user/tweets`, {})
  }

  getPosts(id: number): Promise<Tweet[]> {
    return this.get(`${this.baseUrl}/user/${id}/tweets`, {})
  }

  getPost(id: number): Promise<Tweet> {
    return this.get(`${this.baseUrl}/tweet/${id}`, {})
  }

  likePost(id: number): Promise<any> {
    return this.post(`${this.baseUrl}/tweet/${id}/like`, {})
  }

  follow(id: number): Promise<any> {
    return this.post(`${this.baseUrl}/user/${id}/follow`, {})
  }

  createUserIfNotExists(): Promise<User> {
    return this.post(`${this.baseUrl}/user`, {})
  }
}
