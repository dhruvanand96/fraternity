import { Injectable } from '@angular/core';
import { CanActivate, Router, UrlTree } from '@angular/router';
import { TokenStorageService } from '../services/token-storage.service';

@Injectable({ providedIn: 'root' })
export class AuthGuard implements CanActivate {
  constructor(
    private readonly tokens: TokenStorageService,
    private readonly router: Router
  ) {}

  canActivate(): boolean | UrlTree {
    return this.tokens.isAuthenticated
      ? true
      : this.router.createUrlTree(['/login']);
  }
}
