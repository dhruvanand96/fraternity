import { environment } from './../../environments/environment';
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Store } from '@ngrx/store';
import { map } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class UserService {
  constructor(private http: HttpClient,) { }


  getUsers(): Observable<any> {
    return this.http.get(environment.serverUrl+'/show-users')
  }

  deleteUser(params: any): Observable<any> {
    return this.http.post(environment.serverUrl+'/delete', params);
  }

}
