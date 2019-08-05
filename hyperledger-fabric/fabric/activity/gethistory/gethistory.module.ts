
import {
   WiServiceContribution
} from "wi-studio/app/contrib/wi-contrib";
import { NgModule } from "@angular/core";
import { CommonModule } from "@angular/common";
import { HttpModule, Http } from "@angular/http";
import { gethistoryHandler } from "./gethistoryHandler";
@NgModule({
  imports: [
    CommonModule,
HttpModule
  ],
  exports: [],
  declarations: [],
  entryComponents: [],
  providers: [
    { provide: WiServiceContribution, useClass: gethistoryHandler}
  ],
  bootstrap: []
})
export default class gethistoryModule {

}
