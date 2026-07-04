import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
import { LoginComponent } from './login/login.component';
import { RegisterComponent } from './register/register.component';
import { ReactiveFormsModule } from '@angular/forms';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import {MatInputModule} from '@angular/material/input';
import { HomeComponent } from './home/home.component'
import {MatTableModule} from '@angular/material/table';
import { AuthInterceptor } from './core/interceptors/auth.interceptor';
import {MatProgressSpinnerModule} from '@angular/material/progress-spinner';
import { ToastrModule } from 'ngx-toastr';
import { StoreModule } from '@ngrx/store';
import { userReducer } from './store/users/user.reducer';
import { EffectsModule } from '@ngrx/effects';
import { UserEffects } from './store/users/user.effects';


@NgModule({
  declarations: [
    AppComponent,
    LoginComponent,
    RegisterComponent,
    HomeComponent
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    AppRoutingModule,
    ReactiveFormsModule,
    BrowserAnimationsModule,
    MatInputModule,
    MatTableModule,
    MatProgressSpinnerModule,
    StoreModule.forRoot({ user: userReducer }),
    EffectsModule.forRoot([UserEffects]),
    ToastrModule.forRoot()

  ],
  providers: [{
    provide: HTTP_INTERCEPTORS,
    useClass: AuthInterceptor,
    multi: true
  }],
  bootstrap: [AppComponent]
})
export class AppModule { }
