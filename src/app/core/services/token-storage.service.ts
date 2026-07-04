import { Injectable } from '@angular/core';

const TOKEN_KEY = 'login-token';
const USER_KEY = 'user';

@Injectable({ providedIn: 'root' })
export class TokenStorageService {
  get token(): string | null {
    return localStorage.getItem(TOKEN_KEY);
  }

  save(token: string, user: unknown): void {
    localStorage.setItem(TOKEN_KEY, token);
    localStorage.setItem(USER_KEY, JSON.stringify(user));
  }

  clear(): void {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
  }

  get isAuthenticated(): boolean {
    return Boolean(this.token);
  }
}
