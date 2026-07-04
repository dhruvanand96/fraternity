import {
  HttpErrorResponse,
  HttpEvent,
  HttpHandler,
  HttpInterceptor,
  HttpRequest
} from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { ToastrService } from 'ngx-toastr';
import { Observable, throwError } from 'rxjs';
import { catchError, finalize } from 'rxjs/operators';
import { LoadingService } from '../services/loading.service';
import { TokenStorageService } from '../services/token-storage.service';

@Injectable()
export class AuthInterceptor implements HttpInterceptor {
  private requestSequence = 0;

  constructor(
    private readonly router: Router,
    private readonly loading: LoadingService,
    private readonly tokens: TokenStorageService,
    private readonly toastr: ToastrService
  ) {}

  intercept(request: HttpRequest<unknown>, next: HttpHandler): Observable<HttpEvent<unknown>> {
    const requestId = `${request.method}:${request.url}:${this.requestSequence++}`;
    const token = this.tokens.token;
    const authenticatedRequest = token
      ? request.clone({ setHeaders: { Authorization: `Bearer ${token}` } })
      : request;

    this.loading.start(requestId);
    return next.handle(authenticatedRequest).pipe(
      catchError((error: HttpErrorResponse) => {
        if (error.status === 401) {
          this.tokens.clear();
          void this.router.navigate(['/login']);
          this.toastr.error('Your session is invalid or has expired.');
        } else {
          this.toastr.error(error.error?.error || 'The request could not be completed.');
        }
        return throwError(error);
      }),
      finalize(() => this.loading.stop(requestId))
    );
  }
}
