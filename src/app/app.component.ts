import { Component } from '@angular/core';
import { hinos } from './hinos'

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'harpa';
  hinos = {};
  hino = null
  nome = ""
  letra = ""
  selected = 0;

  constructor() {
    hinos.forEach(hino => {
      this.hinos[hino.id] = hino;
    });
  }

  public add(val: number): void {
    this.selected = this.selected * 10 + val;
    this.nome = this.hinos[this.selected].nome.toUpperCase();
  }

  public clear(val: Number): void {
    this.selected = 0;
    this.hino = this.hinos[this.selected]
    this.nome = "";
  }

  public show(): void {
    this.hino = this.hinos[this.selected]
    this.nome = this.hino.nome.toUpperCase()
    this.letra = this.hino.letra
  }
}
