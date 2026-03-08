import { ComponentFixture, TestBed } from '@angular/core/testing';
import { of } from 'rxjs';
import { InventoryDashboard } from './inventory-dashboard';
import { Api } from '../api';

describe('InventoryDashboard', () => {
  let component: InventoryDashboard;
  let fixture: ComponentFixture<InventoryDashboard>;

  beforeEach(async () => {
    const mockApi = { getResources: () => of([]) };
    await TestBed.configureTestingModule({
      imports: [InventoryDashboard],
      providers: [{ provide: Api, useValue: mockApi }]
    }).compileComponents();

    fixture = TestBed.createComponent(InventoryDashboard);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});