import { WiServiceContribution } from "wi-studio/app/contrib/wi-contrib";
import { NgModule } from "@angular/core";
import { CommonModule } from "@angular/common";
import { HttpModule, Http } from "@angular/http";
import { transactionHandler } from "./transactionHandler";
@NgModule({
  imports: [CommonModule, HttpModule],
  exports: [],
  declarations: [],
  entryComponents: [],
  providers: [{ provide: WiServiceContribution, useClass: transactionHandler }],
  bootstrap: [],
})
export default class transactionModule {}
