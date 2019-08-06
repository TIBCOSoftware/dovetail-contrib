
import {
   WiServiceContribution
} from "wi-studio/app/contrib/wi-contrib";
import { NgModule } from "@angular/core";
import { CommonModule } from "@angular/common";
import { HttpModule, Http } from "@angular/http";
import { getbycompositekeyHandler } from "./getbycompositekeyHandler";
@NgModule({
  imports: [
    CommonModule,
HttpModule
  ],
  exports: [],
  declarations: [],
  entryComponents: [],
  providers: [
    { provide: WiServiceContribution, useClass: getbycompositekeyHandler}
  ],
  bootstrap: []
})
export default class getbycompositekeyModule {

}
