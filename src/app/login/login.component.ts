import { AuthService } from './../_services/auth.service';
import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { LoadingService } from '../_services/loading.service';
import { ToastrService } from 'ngx-toastr';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {

  constructor(private authService: AuthService,
              private router: Router,
              private _loading: LoadingService,
              private toastr: ToastrService) { }


  loginForm = new FormGroup({
    email: new FormControl('',Validators.pattern("^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}$")),
    password: new FormControl('',Validators.required),
  });

  ngOnInit(): void {
  }


  onSubmit() {
    // TODO: Use EventEmitter with form value
    this.authService.login(this.loginForm.value).subscribe( (res: any) => {
      if ( res["access-token"] != "" && res.error == ""){
        this.toastr.success("Login Successful..!", "")
        localStorage.setItem('login-token', res["access-token"])
        localStorage.setItem('user', res["user"])
        this.router.navigate(['home']);
      }
    })
  }

}
