import { ComponentFixture, TestBed } from '@angular/core/testing';
import { of } from 'rxjs';
import { ProductionChainVisualization } from './production-chain-visualization';
import { Api } from '../api';

describe('ProductionChainVisualization', () => {
  let component: ProductionChainVisualization;
  let fixture: ComponentFixture<ProductionChainVisualization>;

  beforeEach(async () => {
    const mockApi = { getProductionChains: () => of([]) };
    await TestBed.configureTestingModule({
      imports: [ProductionChainVisualization],
      providers: [{ provide: Api, useValue: mockApi }]
    }).compileComponents();

    fixture = TestBed.createComponent(ProductionChainVisualization);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});