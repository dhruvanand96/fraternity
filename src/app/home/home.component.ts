import { UserService } from './../_services/user.service';
import { Component, OnInit } from '@angular/core';
import { AuthService } from '../_services/auth.service';
import { Router } from '@angular/router';
import { Store, select } from '@ngrx/store';
import { removeUserById } from '../store/selector/user.selectors';
import { invokeDeleteAPI } from '../store/action/user.actions';


export interface Element {
  Name: string;
  email: string;
  id: number;
  password: string;
}

const ELEMENT_DATA: Element[] = [
];


@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})

export class HomeComponent implements OnInit {


  selectedAlbumId = -1;
  allUsers$ = this.store.pipe(
    select(removeUserById(this.selectedAlbumId))
  );

  displayedColumns: string[] = ['id','name','email', 'actions' ];
  dataSource = ELEMENT_DATA;

  constructor( private _UserService : UserService,
               private authService: AuthService,
               private router: Router,
               private store: Store<{ user: any[] }>,) {
   }


  ngOnInit(): void {
    this._UserService.getUsers().subscribe( (res: any) => {

      if (res.user_array != null){
        this.dataSource =  res.user_array
      }
    })
  }


  submit(){
    this.authService.logout().subscribe( (res: any) => {
        this.router.navigate(['login']);
    })
  }

  deleteUser(user: any){
    this.store.dispatch(invokeDeleteAPI(user));
  // this._UserService.deleteUser(user).subscribe( (res: any) =>{})
  this._UserService.getUsers().subscribe( (res: any) => {

    if (res.user_array != null){
      this.dataSource =  res.user_array
    }
  })
  }

}
