import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Resource {
  id: string;
  name: string;
  quantity: number;
  quality: number;
  category: string;
}

export interface ProductionChain {
  id: string;
  name: string;
  inputs: Resource[];
  outputs: Resource[];
}

export interface Season {
  id: string;
  name: string;
  weather: string;
  agriculturalCalendar: string[];
}

@Injectable({
  providedIn: 'root',
})
export class Api {
  private baseUrl = 'http://localhost:8080/api';

  constructor(private http: HttpClient) {}

  // Resource endpoints
  getResources(): Observable<Resource[]> {
    return this.http.get<Resource[]>(`${this.baseUrl}/resources`);
  }

  getResource(id: string): Observable<Resource> {
    return this.http.get<Resource>(`${this.baseUrl}/resources/${id}`);
  }

  // Production chain endpoints
  getProductionChains(): Observable<ProductionChain[]> {
    return this.http.get<ProductionChain[]>(`${this.baseUrl}/production-chains`);
  }

  // Seasonal data
  getCurrentSeason(): Observable<Season> {
    return this.http.get<Season>(`${this.baseUrl}/season/current`);
  }

  // Quality tiers
  getQualityTiers(): Observable<{ tier: number; label: string; color: string }[]> {
    return this.http.get<{ tier: number; label: string; color: string }[]>(`${this.baseUrl}/quality-tiers`);
  }
}