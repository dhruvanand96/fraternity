import { Component, OnInit } from '@angular/core';
import { delay } from 'rxjs/operators';
import { DataServiceService } from './data-service.service';
import { LoadingService } from './_services/loading.service';


@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})


export class AppComponent implements OnInit {

  title: any;
  loading: boolean = false;

  constructor(
    private _loading: LoadingService
  ){ }

  ngOnInit() {
    this.listenToLoading();
  }

  /**
   * Listen to the loadingSub property in the LoadingService class. This drives the
   * display of the loading spinner.
   */
 listenToLoading(): void {
  this._loading.loadingSub
    .pipe(delay(0)) // This prevents a ExpressionChangedAfterItHasBeenCheckedError for subsequent requests
    .subscribe((loading) => {
      this.loading = loading;
    });
}

}
