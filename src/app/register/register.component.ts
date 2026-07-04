import { Router } from '@angular/router';
import { Component, OnInit } from '@angular/core';
import { FormGroup, FormControl, Validators } from '@angular/forms';
import { AuthService, Registration } from '../core/services/auth.service';
import { TokenStorageService } from '../core/services/token-storage.service';
import { ToastrService } from 'ngx-toastr';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.css']
})
export class RegisterComponent implements OnInit {


  registerForm = new FormGroup({

    name: new FormControl('', Validators.required),
    email: new FormControl('', [Validators.required, Validators.email]),
    password: new FormControl('', Validators.required),
    confirmPassword: new FormControl('', Validators.required)

  });
  constructor(
    private readonly router: Router,
    private readonly authService: AuthService,
    private readonly tokens: TokenStorageService,
    private readonly toastr: ToastrService
  ) {}

  ngOnInit(): void {
  }

  onSubmit(): void {
    if (this.registerForm.invalid) {
      return;
    }
    const values = this.registerForm.value as Registration & { confirmPassword: string };
    if (values.password !== values.confirmPassword) {
      this.toastr.error('Passwords do not match.');
      return;
    }
    const registration: Registration = {
      name: values.name,
      email: values.email,
      password: values.password
    };
    this.authService.register(registration).subscribe(() => {
      this.authService.login(registration).subscribe(response => {
        this.tokens.save(response['access-token'], response.user);
        this.toastr.success('User successfully created.');
        void this.router.navigate(['/home']);
      });
    });
  }
}
