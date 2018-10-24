import { HttpModule, Http } from "@angular/http";
import { NgModule } from "@angular/core";
import { KVLedgerActivityService, Field1Provider, Field1ValidationProvider } from "./activity";
import { WiServiceContribution} from "wi-studio/app/contrib/wi-contrib";


@NgModule({
  imports: [
    HttpModule,
  ],
  exports: [

  ],
  declarations: [

  ],
  entryComponents: [

  ],
  providers: [
    {
       provide: WiServiceContribution,
       useClass: KVLedgerActivityService
     },
     {
       provide: Field1Provider,
       useClass: Field1Provider
     },
     {
       provide: Field1ValidationProvider,
       useClass: Field1ValidationProvider
     }
  ],
  bootstrap: []
})
export default class KVLedgerContribModule {

}
