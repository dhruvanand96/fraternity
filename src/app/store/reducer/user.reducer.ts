import { Action, createReducer, on } from '@ngrx/store';
import { invokeDeleteAPI, SuccessGetUserAction, usersList } from '../action/user.actions';


export const userFeatureKey = 'user';

export interface State {

}

export const initialState: ReadonlyArray<any> = [];

const _userReducer = createReducer(
  initialState,
  on(invokeDeleteAPI, (state, { user }) => {
    return [user];
  }),
  on(SuccessGetUserAction, (state, {
    payload
 }) => {
    return {
       ...state,
       users: payload,
       isLoaded: true
    };
 }),
);

export function userReducer(state: any, action: any) {
  return _userReducer(state, action);
}