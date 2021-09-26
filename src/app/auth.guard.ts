import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, CanActivate, CanActivateChild, Router, RouterStateSnapshot, UrlTree } from '@angular/router';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class AuthGuard implements CanActivate, CanActivateChild {

  constructor(private _Router: Router){}


  canActivateChild(): boolean {
    if (localStorage.getItem('login-token') != '') {
      return true;
    }
    else {
      localStorage.removeItem('login-token');
      this._Router.navigate(['/login']);
      return false;
    }
  }

  CheckToken() {
    return !!localStorage.getItem('login-token');
  }

  canActivate(): boolean {
    if (this.CheckToken()) {
      return true;
    }
    else {
      this._Router.navigate(['/login']);
      return false;
    }
  }
}
