import { ComponentFixture, TestBed } from '@angular/core/testing';
import { of } from 'rxjs';
import { SeasonalIndicators } from './seasonal-indicators';
import { Api } from '../api';

describe('SeasonalIndicators', () => {
  let component: SeasonalIndicators;
  let fixture: ComponentFixture<SeasonalIndicators>;

  beforeEach(async () => {
    const mockApi = { getCurrentSeason: () => of({ id: 'spring', name: 'Spring', weather: 'Rainy', agriculturalCalendar: [] }) };
    await TestBed.configureTestingModule({
      imports: [SeasonalIndicators],
      providers: [{ provide: Api, useValue: mockApi }]
    }).compileComponents();

    fixture = TestBed.createComponent(SeasonalIndicators);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});