import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { Economy } from './economy/economy/economy';

const routes: Routes = [
  { path: '', component: Economy }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class EconomyRoutingModule {}