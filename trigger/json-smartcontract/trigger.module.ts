import { HttpModule } from '@angular/http';
import { NgModule } from "@angular/core";
import { JsonSmartContractTriggerService } from "./trigger";
import { WiServiceContribution } from "wi-studio/app/contrib/wi-contrib";

@NgModule({
  imports: [
    HttpModule
  ],
  providers: [
    {
       provide: WiServiceContribution,
       useClass: JsonSmartContractTriggerService
    }
  ]
})
export default class JsonSmartContractContribModule {

}
