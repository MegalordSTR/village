import { Component } from '@angular/core';
import { ProductionChainVisualization } from '../production-chain-visualization/production-chain-visualization';
import { InventoryDashboard } from '../inventory-dashboard/inventory-dashboard';
import { SeasonalIndicators } from '../seasonal-indicators/seasonal-indicators';
import { QualityDisplay } from '../quality-display/quality-display';

@Component({
  selector: 'app-economy',
  imports: [ProductionChainVisualization, InventoryDashboard, SeasonalIndicators, QualityDisplay],
  templateUrl: './economy.html',
  styleUrl: './economy.css',
})
export class Economy {}