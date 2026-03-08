import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Api, Season } from '../api';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-seasonal-indicators',
  imports: [CommonModule],
  templateUrl: './seasonal-indicators.html',
  styleUrl: './seasonal-indicators.css',
})
export class SeasonalIndicators implements OnInit {
  season$: Observable<Season>;

  constructor(private api: Api) {
    this.season$ = this.api.getCurrentSeason();
  }

  ngOnInit(): void {}
}