import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Api } from '../api';
import { Observable } from 'rxjs';

interface QualityTier {
  tier: number;
  label: string;
  color: string;
}

@Component({
  selector: 'app-quality-display',
  imports: [CommonModule],
  templateUrl: './quality-display.html',
  styleUrl: './quality-display.css',
})
export class QualityDisplay implements OnInit {
  tiers$: Observable<QualityTier[]>;

  constructor(private api: Api) {
    this.tiers$ = this.api.getQualityTiers();
  }

  ngOnInit(): void {}
}