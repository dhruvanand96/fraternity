import { Action, createAction, props } from '@ngrx/store';


export const usersList = createAction(
  "[USER] Remove",
  props<{allusers:any[]}>()
);


export const invokeDeleteAPI = createAction(
'[Delete API] Invoke API',
props<{user:any}>()
);

export const SuccessGetUserAction = createAction(
  '[User] - Success',
  props<{ payload: any[] }>()
);

export const ErrorUserAction = createAction('[User] - Error', props<Error>());