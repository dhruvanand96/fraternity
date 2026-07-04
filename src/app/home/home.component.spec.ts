import { NO_ERRORS_SCHEMA } from '@angular/core';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { Store } from '@ngrx/store';
import { of } from 'rxjs';
import { AuthService } from '../core/services/auth.service';
import { HomeComponent } from './home.component';

describe('HomeComponent', () => {
  let fixture: ComponentFixture<HomeComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [RouterTestingModule],
      declarations: [HomeComponent],
      providers: [
        { provide: AuthService, useValue: { logout: () => of({}) } },
        { provide: Store, useValue: { pipe: () => of([]), dispatch: () => undefined } }
      ],
      schemas: [NO_ERRORS_SCHEMA]
    }).compileComponents();
    fixture = TestBed.createComponent(HomeComponent);
  });

  it('creates and requests users on initialization', () => {
    const store = TestBed.inject(Store);
    spyOn(store, 'dispatch');
    fixture.detectChanges();
    expect(fixture.componentInstance).toBeTruthy();
    expect(store.dispatch).toHaveBeenCalled();
  });
});
