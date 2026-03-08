import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Api, ProductionChain } from '../api';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-production-chain-visualization',
  imports: [CommonModule],
  templateUrl: './production-chain-visualization.html',
  styleUrl: './production-chain-visualization.css',
})
export class ProductionChainVisualization implements OnInit {
  chains$: Observable<ProductionChain[]>;

  constructor(private api: Api) {
    this.chains$ = this.api.getProductionChains();
  }

  ngOnInit(): void {}
}