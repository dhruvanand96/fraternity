import { Router } from '@angular/router';
import { HTTP_INTERCEPTORS, HttpEvent, HttpResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { HttpInterceptor, HttpHandler, HttpRequest } from '@angular/common/http';
import { Observable } from 'rxjs';
import { tap } from 'rxjs/internal/operators/tap';
import { LoadingService } from '../loading.service';
import { map } from 'rxjs/internal/operators/map';
import { ToastrService } from 'ngx-toastr';

@Injectable()
export class AuthInterceptor implements HttpInterceptor {

  headers: any
  constructor(private _router: Router,
              private _loading: LoadingService,
              private toastr: ToastrService) { }

  intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    this._loading.setLoading(true, request.url);
  let token = localStorage.getItem("login-token");
  this.headers = {
    setHeaders: {
        "Authorization": "BEARER "+token
    }
}
 const clonedreq = request.clone(this.headers);
  return next.handle(clonedreq).pipe(tap(
    succ => {
        if (request.method != 'GET' && request.url.indexOf('logout') == -1 && request.url.indexOf('login') == -1) {
            this._loading.setLoading(false, request.url);
            if (succ instanceof HttpResponse) {
        }
    }
        if (succ instanceof HttpResponse) {
            this._loading.setLoading(false, request.url);
        }
    },
    err => {
        this._loading.setLoading(false, request.url);
        if (err.status === 401) {
            this.toastr.error("Invalid credentials", "")
           this.clearData();
        }
        else if (err.status === 404) {
            this.toastr.error("Error Please try again later..!", "")
            const validationError = err.error;
        }
        else if (err.status === 400) {
            this.toastr.error("Error Please try again later..!", "")
            const validationError = err.error;
        }
    }
  ))
  .pipe(map<HttpEvent<any>, any>((evt: HttpEvent<any>) => {
    if (evt instanceof HttpResponse) {
      this._loading.setLoading(false, request.url);
    }
    return evt;
  }));
}

clearData() {;
    localStorage.clear();
    this._router.navigate(['/login']);
}
}



export const authInterceptorProviders = [
  { provide: HTTP_INTERCEPTORS, useClass: AuthInterceptor, multi: true }
];