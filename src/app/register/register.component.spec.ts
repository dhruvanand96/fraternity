import { NO_ERRORS_SCHEMA } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { of } from 'rxjs';
import { ToastrService } from 'ngx-toastr';
import { AuthService } from '../core/services/auth.service';
import { RegisterComponent } from './register.component';

describe('RegisterComponent', () => {
  let fixture: ComponentFixture<RegisterComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ReactiveFormsModule, RouterTestingModule],
      declarations: [RegisterComponent],
      providers: [
        { provide: AuthService, useValue: { register: () => of({}), login: () => of({}) } },
        { provide: ToastrService, useValue: { success: () => undefined, error: () => undefined } }
      ],
      schemas: [NO_ERRORS_SCHEMA]
    }).compileComponents();
    fixture = TestBed.createComponent(RegisterComponent);
  });

  it('requires all registration fields', () => {
    const component = fixture.componentInstance;
    expect(component).toBeTruthy();
    expect(component.registerForm.valid).toBeFalse();
  });
});
