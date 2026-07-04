import { deleteUserSuccess, loadUsersSuccess } from './user.actions';
import { initialState, userReducer } from './user.reducer';

describe('userReducer', () => {
  it('stores loaded users and removes a deleted user', () => {
    const loaded = userReducer(initialState, loadUsersSuccess({ users: [
      { id: 1, name: 'One', email: 'one@example.com' },
      { id: 2, name: 'Two', email: 'two@example.com' }
    ] }));
    const result = userReducer(loaded, deleteUserSuccess({ id: 1 }));
    expect(result.users.map(user => user.id)).toEqual([2]);
  });
});
