import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Api, Resource } from '../api';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-inventory-dashboard',
  imports: [CommonModule],
  templateUrl: './inventory-dashboard.html',
  styleUrl: './inventory-dashboard.css',
})
export class InventoryDashboard implements OnInit {
  resources$: Observable<Resource[]>;

  constructor(private api: Api) {
    this.resources$ = this.api.getResources();
  }

  ngOnInit(): void {}
}