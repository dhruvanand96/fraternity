import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class LoadingService {
  private readonly pendingRequests = new Set<string>();
  private readonly loadingSubject = new BehaviorSubject<boolean>(false);

  readonly loading$: Observable<boolean> = this.loadingSubject.asObservable();

  start(requestId: string): void {
    this.pendingRequests.add(requestId);
    this.loadingSubject.next(true);
  }

  stop(requestId: string): void {
    this.pendingRequests.delete(requestId);
    this.loadingSubject.next(this.pendingRequests.size > 0);
  }
}
