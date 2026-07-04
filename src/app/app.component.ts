import { Component } from '@angular/core';
import { Observable } from 'rxjs';
import { LoadingService } from './core/services/loading.service';


@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})


export class AppComponent {
  readonly title = 'fraternity';
  readonly loading$: Observable<boolean>;

  constructor(loading: LoadingService) {
    this.loading$ = loading.loading$;
  }
}
