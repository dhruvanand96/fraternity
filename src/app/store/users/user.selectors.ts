import { createFeatureSelector, createSelector } from '@ngrx/store';
import { UserState } from './user.reducer';

export const selectUserState = createFeatureSelector<UserState>('user');
export const selectUsers = createSelector(selectUserState, state => state.users);
export const selectUsersLoading = createSelector(selectUserState, state => state.loading);
export const selectUsersError = createSelector(selectUserState, state => state.error);
