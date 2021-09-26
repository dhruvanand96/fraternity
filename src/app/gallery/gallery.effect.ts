import { Injectable } from '@angular/core';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { of } from 'rxjs';
import { catchError, map, mergeMap } from 'rxjs/operators';
import { UserService } from '../_services/user.service';
import * as UserActions from '../store/action/user.actions';
@Injectable()
export class GalleryEffect {
  constructor(
    private actions$: Actions,
    private userService: UserService
  ) {}

  loadusers$ = createEffect(() =>
    this.actions$.pipe(
      ofType('[Delete API] Invoke API'),
      mergeMap((id) =>
        this.userService
          .deleteUser(id)
          .pipe(
              map((data : any[]) => {
                debugger
                return UserActions.SuccessGetUserAction({ payload: data })
              }),
              catchError((error: Error) => {
               return of(UserActions.ErrorUserAction(error));
          })
        ))
      )
  );
}