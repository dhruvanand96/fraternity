import { Action, createReducer, on } from '@ngrx/store';
import { User } from '../../models/user';
import * as UserActions from './user.actions';

export interface UserState {
  users: User[];
  loading: boolean;
  error: string | null;
}

export const initialState: UserState = {
  users: [],
  loading: false,
  error: null
};

const reducer = createReducer(
  initialState,
  on(UserActions.loadUsers, UserActions.deleteUser, state => ({
    ...state,
    loading: true,
    error: null
  })),
  on(UserActions.loadUsersSuccess, (state, { users }) => ({
    ...state,
    users,
    loading: false
  })),
  on(UserActions.deleteUserSuccess, (state, { id }) => ({
    ...state,
    users: state.users.filter(user => user.id !== id),
    loading: false
  })),
  on(UserActions.loadUsersFailure, UserActions.deleteUserFailure, (state, { error }) => ({
    ...state,
    loading: false,
    error
  }))
);

export function userReducer(state: UserState | undefined, action: Action): UserState {
  return reducer(state, action);
}
