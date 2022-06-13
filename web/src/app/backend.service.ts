import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import config from './app.config';

@Injectable({
  providedIn: 'root',
})
export class BackendService {
  private readonly host: string;

  constructor(private http: HttpClient) {
    this.host = config.backend.host;
  }

  getProfile(): Observable<Profile> {
    const url = new URL('/api/v1/profile', this.host);
    return this.http.get<Profile>(url.href);
  }
}

export interface Profile {
  email: string;
}
