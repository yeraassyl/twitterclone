export interface IAuthResponse {
  auth_token: string
  email: string
}

export interface User {
  id: number
  email: string
  followers: number[]
  following: number[]
}

export interface UserSimple {
  id: number
  email: string
}

export class UserSimplImpl implements UserSimple{
  email: string;
  id: number;

  constructor(email: string, id: number) {
    this.email = email;
    this.id = id;
  }
}

export interface Tweet {
  id: number
  timePosted: string
  content: string
  likes: number
}
