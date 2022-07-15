import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import config from './app.config';

@Injectable({
  providedIn: 'root',
})
export class ValidatorService {
  private readonly host: string;

  constructor(private http: HttpClient) {
    this.host = config.backend.validator;
  }

  validateTag(validate: ValidateRequest): Observable<ValidateResponse> {
    if (this.host === '') {
      return new Observable<ValidateResponse>();
    }

    const url = new URL(`/api/v1/validate/tags`, this.host);

    return this.http.post<ValidateResponse>(url.href, validate);
  }
}

export interface ValidateRequest {
  data: { [key: string]: string };
}

export interface ValidateResponse {
  data: { [key: string]: string };
}
