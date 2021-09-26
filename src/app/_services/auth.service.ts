import { environment } from './../../environments/environment';
import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';


const httpOptions = {
  headers: new HttpHeaders({
    'Content-Type': 'multipart/form-data'
  })
};

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  constructor(private http: HttpClient) { }

  login(params: any): Observable<any> {
    return this.http.post<any>(environment.serverUrl+'/login', params ,httpOptions );
  }

  logout(): Observable<any> {
    return this.http.get<any>(environment.serverUrl+'/logout' ,httpOptions );
  }

  register(params: any): Observable<any> {
    return this.http.post(environment.serverUrl+'/create-user', params, httpOptions);
  }
}