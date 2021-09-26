import { createFeatureSelector, createSelector } from '@ngrx/store';
import { AppState } from '../app.state';

export const userSelector =(state: AppState) => state.user;



export const removeUserById = ( id:number) => createSelector(
    userSelector,
    (user:any[]) => {
        if(id == -1){
            return user;
        }
        return user.filter(_ => _.id == id);
    }
)