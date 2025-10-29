import { Component } from '@angular/core';
import { ApiService } from '../../services/api.service';

@Component({
  selector: 'app-scraper',
  templateUrl: './scraper.component.html',
  styleUrls: ['./scraper.component.css']
})
export class ScraperComponent {
  url = '';
  question = '';
  cleaned = '';
  answer = '';
  loading = false;
  model = 'gemma2:2b';

  availableModels = [
    'gemma2:2b',
    'gemma3:12b',
    'llama3.2',
  ];

  constructor(private api: ApiService) {}

  onScrape() {
    if (!this.url) return;
    this.loading = true;
    this.cleaned = '';
    this.answer = '';
    this.api.scrape(this.url).subscribe({
      next: (res) => {
        this.cleaned = res?.cleaned || '';
        this.loading = false;
      },
      error: (err) => {
        alert('Scrape failed: ' + (err?.error?.error || err.message));
        this.loading = false;
      }
    });
  }

  onParse() {
    if (!this.question) return;
    this.loading = true;
    this.answer = '';
    const dom = this.cleaned?.trim().length ? this.cleaned : undefined;
    const url = !dom && this.url ? this.url : undefined;
    this.api.parse(this.question, dom, url, this.model).subscribe({
      next: (res) => {
        this.answer = res?.answer || '';
        this.loading = false;
      },
      error: (err) => {
        alert('Parse failed: ' + (err?.error?.error || err.message));
        this.loading = false;
      }
    });
  }
}
