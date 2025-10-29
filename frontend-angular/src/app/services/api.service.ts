import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class ApiService {
  constructor(private http: HttpClient) {}

  scrape(url: string): Observable<any> {
    return this.http.post('/api/scrape', { url });
  }

  parse(question: string, domContent?: string, url?: string, model?: string): Observable<any> {
    return this.http.post('/api/parse', { question, domContent, url, model });
  }
}
