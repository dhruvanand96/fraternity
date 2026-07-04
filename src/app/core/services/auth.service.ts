import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { User } from '../../models/user';

export interface Credentials {
  email: string;
  password: string;
}

export interface Registration extends Credentials {
  name: string;
}

export interface LoginResponse {
  'access-token': string;
  user: User;
  error: string;
}

@Injectable({ providedIn: 'root' })
export class AuthService {
  constructor(private readonly http: HttpClient) {}

  login(credentials: Credentials): Observable<LoginResponse> {
    return this.http.post<LoginResponse>(`${environment.serverUrl}/login`, credentials);
  }

  register(registration: Registration): Observable<{ success: string; error: string }> {
    return this.http.post<{ success: string; error: string }>(
      `${environment.serverUrl}/create-user`,
      registration
    );
  }

  logout(): Observable<{ success: string; error: string }> {
    return this.http.get<{ success: string; error: string }>(`${environment.serverUrl}/logout`);
  }
}
