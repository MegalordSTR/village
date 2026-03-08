import { ComponentFixture, TestBed } from '@angular/core/testing';
import { of } from 'rxjs';
import { Economy } from './economy';
import { Api } from '../api';

describe('Economy', () => {
  let component: Economy;
  let fixture: ComponentFixture<Economy>;

  beforeEach(async () => {
    const mockApi = {
      getProductionChains: () => of([]),
      getResources: () => of([]),
      getCurrentSeason: () => of({ id: 'spring', name: 'Spring', weather: 'Rainy', agriculturalCalendar: [] }),
      getQualityTiers: () => of([])
    };
    await TestBed.configureTestingModule({
      imports: [Economy],
      providers: [{ provide: Api, useValue: mockApi }]
    }).compileComponents();

    fixture = TestBed.createComponent(Economy);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
