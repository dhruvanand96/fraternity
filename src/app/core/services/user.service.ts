import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { User } from '../../models/user';

export interface UsersResponse {
  user_array: User[];
  success: string;
  error: string;
}

@Injectable({ providedIn: 'root' })
export class UserService {
  constructor(private readonly http: HttpClient) {}

  list(): Observable<UsersResponse> {
    return this.http.get<UsersResponse>(`${environment.serverUrl}/show-users`);
  }

  delete(id: number): Observable<{ success: string; error: string }> {
    return this.http.post<{ success: string; error: string }>(
      `${environment.serverUrl}/delete`,
      { id }
    );
  }
}
