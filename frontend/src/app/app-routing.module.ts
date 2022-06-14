import {RouterModule, Routes} from "@angular/router";
import {PostsComponent} from "./posts/posts.component";
import {UsersComponent} from "./users/users.component";
import {UserComponent} from "./user/user.component";
import {NgModule} from "@angular/core";
import {PostComponent} from "./post/post.component";


const routes: Routes = [
  {path: 'posts', component: PostsComponent},
  {path: 'users', component: UsersComponent},
  {path: 'user/:id', component: UserComponent},
  {path: 'post/:id', component: PostComponent}
]

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})

export class AppRoutingModule {}
