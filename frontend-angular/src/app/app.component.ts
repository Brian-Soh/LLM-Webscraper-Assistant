import { Component } from '@angular/core';

@Component({
  selector: 'app-root',
  template: `
    <div class="container">
      <h1>LLM Webscraper Assistant</h1>
      <app-scraper></app-scraper>
    </div>
  `,
  styles: [`
    .container {
      max-width: 900px;
      margin: 2rem auto;
      font-family: system-ui, -apple-system, Segoe UI, Roboto, Helvetica, Arial;
      text-align: center; /* centers all text inside */
    }

    h1 {
      font-weight: 600;
      margin-bottom: 1.5rem;
    }
  `]
})
export class AppComponent {}
