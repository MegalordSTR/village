import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { Api, Resource, ProductionChain, Season } from './api';

describe('Api', () => {
  let service: Api;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [Api]
    });
    service = TestBed.inject(Api);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  describe('getResources', () => {
    it('should fetch resources', () => {
      const mockResources: Resource[] = [
        { id: '1', name: 'Wood', quantity: 100, quality: 2, category: 'raw' }
      ];

      service.getResources().subscribe(resources => {
        expect(resources).toEqual(mockResources);
      });

      const req = httpMock.expectOne(`${service['baseUrl']}/resources`);
      expect(req.request.method).toBe('GET');
      req.flush(mockResources);
    });
  });

  describe('getResource', () => {
    it('should fetch a single resource', () => {
      const mockResource: Resource = { id: '1', name: 'Wood', quantity: 100, quality: 2, category: 'raw' };

      service.getResource('1').subscribe(resource => {
        expect(resource).toEqual(mockResource);
      });

      const req = httpMock.expectOne(`${service['baseUrl']}/resources/1`);
      expect(req.request.method).toBe('GET');
      req.flush(mockResource);
    });
  });

  describe('getProductionChains', () => {
    it('should fetch production chains', () => {
      const mockChains: ProductionChain[] = [
        { id: '1', name: 'Lumber', inputs: [], outputs: [] }
      ];

      service.getProductionChains().subscribe(chains => {
        expect(chains).toEqual(mockChains);
      });

      const req = httpMock.expectOne(`${service['baseUrl']}/production-chains`);
      expect(req.request.method).toBe('GET');
      req.flush(mockChains);
    });
  });

  describe('getCurrentSeason', () => {
    it('should fetch current season', () => {
      const mockSeason: Season = { id: 'spring', name: 'Spring', weather: 'Rainy', agriculturalCalendar: [] };

      service.getCurrentSeason().subscribe(season => {
        expect(season).toEqual(mockSeason);
      });

      const req = httpMock.expectOne(`${service['baseUrl']}/season/current`);
      expect(req.request.method).toBe('GET');
      req.flush(mockSeason);
    });
  });

  describe('getQualityTiers', () => {
    it('should fetch quality tiers', () => {
      const mockTiers = [{ tier: 1, label: 'Poor', color: 'red' }];

      service.getQualityTiers().subscribe(tiers => {
        expect(tiers).toEqual(mockTiers);
      });

      const req = httpMock.expectOne(`${service['baseUrl']}/quality-tiers`);
      expect(req.request.method).toBe('GET');
      req.flush(mockTiers);
    });
  });
});