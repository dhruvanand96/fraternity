import { Router } from '@angular/router';
import { Component, OnInit } from '@angular/core';
import { FormGroup, FormControl, Validators } from '@angular/forms';
import { AuthService } from '../_services/auth.service';
import { ToastrService } from 'ngx-toastr';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.css']
})
export class RegisterComponent implements OnInit {


  registerForm = new FormGroup({

    name: new FormControl('',Validators.required),
    email: new FormControl('',Validators.pattern("^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}$")),
    password: new FormControl('',Validators.required),
    confirmPassword:  new FormControl('',Validators.required)

  });
  constructor(private router : Router,
              private authService : AuthService,
              private toastr: ToastrService ) { }

  ngOnInit(): void {
  }

  onSubmit() {
    // TODO: Use EventEmitter with form value

    if (this.registerForm.value.password != this.registerForm.value.confirmPassword){
      this.toastr.error("Password do not match", "")
      return
    }
    this.authService.register(this.registerForm.value).subscribe( (res: any) => {
      if (res.error === '' ){
        this.authService.login({"email":this.registerForm.value.email,"password":this.registerForm.value.password}).subscribe(
          (res: any) => {
          if ( res["access-token"] != "" && res.error == ""){
            this.toastr.success("User successfully Created..!", "")
            localStorage.setItem('login-token', res["access-token"])
            this.router.navigate(['home']);
          }else{
            this.toastr.error("Error in Registration", "")
          }
        })
      }
    })
  }

}
