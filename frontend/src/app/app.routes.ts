import { Routes } from '@angular/router';

export const routes: Routes = [
  { path: 'economy', loadChildren: () => import('./economy/economy-module').then(m => m.EconomyModule) },
  { path: '', redirectTo: '/economy', pathMatch: 'full' }
];