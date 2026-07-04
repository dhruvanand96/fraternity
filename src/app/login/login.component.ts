import { AuthService } from '../core/services/auth.service';
import { TokenStorageService } from '../core/services/token-storage.service';
import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { ToastrService } from 'ngx-toastr';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {

  constructor(
    private readonly authService: AuthService,
    private readonly tokens: TokenStorageService,
    private readonly router: Router,
    private readonly toastr: ToastrService
  ) {}


  loginForm = new FormGroup({
    email: new FormControl('', [Validators.required, Validators.email]),
    password: new FormControl('', Validators.required)
  });

  ngOnInit(): void {
  }


  onSubmit(): void {
    if (this.loginForm.invalid) {
      return;
    }
    const credentials = this.loginForm.value as { email: string; password: string };
    this.authService.login(credentials).subscribe(response => {
      this.tokens.save(response['access-token'], response.user);
      this.toastr.success('Login successful.');
      void this.router.navigate(['/home']);
    });
  }
}
