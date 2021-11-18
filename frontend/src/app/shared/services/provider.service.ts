import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {IAuthResponse} from "../models/models";
import * as moment from "moment";

@Injectable({
  providedIn: 'root'
})
export class ProviderService{

  private baseUrl = 'http://localhost:8080'

  constructor(protected http: HttpClient) {
  }

  formatDate(date: Date) {
    return moment(date).format('YYYY-MM-DD');
  }

  private normalBody(body: any): any {
    if (!body) {
      body = {};
    }
    for (const key in body) {
      if (!body.hasOwnProperty(key)) {
        continue;
      }
      if (body[key] instanceof Date) {
        body[key] = this.formatDate(body[key]);
      }
    }
    return body;
  }

  auth(){
    return this.http.get<IAuthResponse>(`http://localhost:8080`).subscribe(resp => {
      console.log(resp);
    })
  }

  listUsers(){
    return this.http.get<any>(`${this.baseUrl}/user`).subscribe(resp => {
      console.log(resp);
    })
  }


}
