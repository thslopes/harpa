import { Component } from '@angular/core';
import { hinos } from './hinos'

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})

export class AppComponent {

  title = 'harpa';
  hinos = Array<any>();
  hino: any;
  nome = "";
  letra = "";
  selected = 0;
  showingAll = false;

  constructor() {
    hinos.forEach(hino => {
      this.hinos[hino.id] = hino;
    });
  }

  public add(val: number): void {
    this.showingAll = false;
    this.selected = this.selected * 10 + val;
    this.nome = this.hinos[this.selected].nome.toUpperCase();
  }

  public clear(): void {
    this.selected = 0;
    this.hino = this.hinos[this.selected];
    this.nome = "";
    this.showingAll = false;
  }

  public showByNumber(n: number): void {
    this.selected = n;
    this.show();
  }

  public show(): void {
    this.hino = this.hinos[this.selected];
    this.nome = this.hino.nome.toUpperCase();
    this.letra = this.hino.letra;
    this.showingAll = false;
  }

  public showAll(): void {
    this.showingAll = true;
  }

  public select(index: number): void {
    this.selected = index;
    this.show();
  }
}
