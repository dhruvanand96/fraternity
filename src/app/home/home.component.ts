import { Component, OnInit } from '@angular/core';
import { AuthService } from '../core/services/auth.service';
import { TokenStorageService } from '../core/services/token-storage.service';
import { Router } from '@angular/router';
import { Store, select } from '@ngrx/store';
import { Observable } from 'rxjs';
import { User } from '../models/user';
import { deleteUser, loadUsers } from '../store/users/user.actions';
import { selectUsers } from '../store/users/user.selectors';


@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})

export class HomeComponent implements OnInit {


  readonly displayedColumns = ['id', 'name', 'email', 'actions'];
  readonly users$: Observable<User[]>;

  constructor(
    private readonly authService: AuthService,
    private readonly tokens: TokenStorageService,
    private readonly router: Router,
    private readonly store: Store
  ) {
    this.users$ = this.store.pipe(select(selectUsers));
  }


  ngOnInit(): void {
    this.store.dispatch(loadUsers());
  }


  logout(): void {
    this.authService.logout().subscribe(() => {
      this.tokens.clear();
      void this.router.navigate(['/login']);
    });
  }

  deleteUser(user: User): void {
    this.store.dispatch(deleteUser({ id: user.id }));
  }
}
