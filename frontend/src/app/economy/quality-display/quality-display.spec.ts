import { ComponentFixture, TestBed } from '@angular/core/testing';
import { of } from 'rxjs';
import { QualityDisplay } from './quality-display';
import { Api } from '../api';

describe('QualityDisplay', () => {
  let component: QualityDisplay;
  let fixture: ComponentFixture<QualityDisplay>;

  beforeEach(async () => {
    const mockApi = { getQualityTiers: () => of([]) };
    await TestBed.configureTestingModule({
      imports: [QualityDisplay],
      providers: [{ provide: Api, useValue: mockApi }]
    }).compileComponents();

    fixture = TestBed.createComponent(QualityDisplay);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});