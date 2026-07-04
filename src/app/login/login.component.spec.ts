import { NO_ERRORS_SCHEMA } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { of } from 'rxjs';
import { ToastrService } from 'ngx-toastr';
import { AuthService } from '../core/services/auth.service';
import { LoginComponent } from './login.component';

describe('LoginComponent', () => {
  let fixture: ComponentFixture<LoginComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ReactiveFormsModule, RouterTestingModule],
      declarations: [LoginComponent],
      providers: [
        { provide: AuthService, useValue: { login: () => of({}) } },
        { provide: ToastrService, useValue: { success: () => undefined, error: () => undefined } }
      ],
      schemas: [NO_ERRORS_SCHEMA]
    }).compileComponents();
    fixture = TestBed.createComponent(LoginComponent);
  });

  it('creates a required, email-valid form', () => {
    const component = fixture.componentInstance;
    expect(component).toBeTruthy();
    expect(component.loginForm.valid).toBeFalse();
    component.loginForm.setValue({ email: 'person@example.com', password: 'secret' });
    expect(component.loginForm.valid).toBeTrue();
  });
});
