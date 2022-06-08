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
  hinosFirstLine = Array<any>();

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

  public search(s: string): void {
    let sn = s.toLowerCase().normalize("NFD").replace(/\p{Diacritic}/gu, "").replace(/[^\w\s]/gu, " ")
    this.hinosFirstLine = Array<any>();
    let hs = this.hinos.slice(1);
    for (let index = 0; index < hs.length; index++) {
      const element = hs[index];
      let i = element.busca.indexOf(sn);
      if (i >= 0) {
        let start = i > 20 ? i - 20 : 0;
        for (; start > 0 && element.busca[start] != ' '; start--) { }
        let end = (start + 50) > element.busca.length ? element.busca.length - 1 : start + 70;
        let pre = start > 0 ? "..." : "";
        let pos = end < element.busca.length - 1 ? "..." : "";
        element.linha = pre + element.busca.slice(start, end) + pos;
        this.hinosFirstLine.push(element);
      }
    }
    this.showingAll = true;
  }

  public showAll(): void {
    this.hinosFirstLine = this.hinos.slice(1);
    this.showingAll = true;
  }

  public select(index: number): void {
    this.selected = index;
    this.show();
  }
}
