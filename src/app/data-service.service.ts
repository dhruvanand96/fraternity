
import { Injectable } from '@angular/core';
import {environment} from '../environments/environment';
import { HttpClient } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError, retry } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class DataServiceService {

 constructor(private http: HttpClient ) { }

 getTitle() {
    return this.http.get(environment.serverUrl+'/dummy-data',{responseType: 'text'})
 }


}
