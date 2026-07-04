import { createAction, props } from '@ngrx/store';
import { User } from '../../models/user';

export const loadUsers = createAction('[Users Page] Load Users');
export const loadUsersSuccess = createAction(
  '[Users API] Load Users Success',
  props<{ users: User[] }>()
);
export const loadUsersFailure = createAction(
  '[Users API] Load Users Failure',
  props<{ error: string }>()
);

export const deleteUser = createAction(
  '[Users Page] Delete User',
  props<{ id: number }>()
);
export const deleteUserSuccess = createAction(
  '[Users API] Delete User Success',
  props<{ id: number }>()
);
export const deleteUserFailure = createAction(
  '[Users API] Delete User Failure',
  props<{ error: string }>()
);
