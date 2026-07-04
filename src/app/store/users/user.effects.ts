import { Injectable } from '@angular/core';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { of } from 'rxjs';
import { catchError, map, switchMap } from 'rxjs/operators';
import { UserService } from '../../core/services/user.service';
import * as UserActions from './user.actions';

@Injectable()
export class UserEffects {
  readonly loadUsers$ = createEffect(() => this.actions$.pipe(
    ofType(UserActions.loadUsers),
    switchMap(() => this.users.list().pipe(
      map(response => UserActions.loadUsersSuccess({ users: response.user_array })),
      catchError(error => of(UserActions.loadUsersFailure({
        error: error.error?.error || 'Could not load users'
      })))
    ))
  ));

  readonly deleteUser$ = createEffect(() => this.actions$.pipe(
    ofType(UserActions.deleteUser),
    switchMap(({ id }) => this.users.delete(id).pipe(
      map(() => UserActions.deleteUserSuccess({ id })),
      catchError(error => of(UserActions.deleteUserFailure({
        error: error.error?.error || 'Could not delete user'
      })))
    ))
  ));

  constructor(
    private readonly actions$: Actions,
    private readonly users: UserService
  ) {}
}
