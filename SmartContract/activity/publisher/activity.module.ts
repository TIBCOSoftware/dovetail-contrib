import { HttpModule } from "@angular/http";
import { NgModule } from "@angular/core";
import { CommonModule } from "@angular/common";

import { EventPublisherActivityContributionHandler} from "./activity";
import { WiServiceContribution} from "wi-studio/app/contrib/wi-contrib";

@NgModule({
  imports: [
    CommonModule,
    HttpModule,
  ],
  exports: [],
  declarations: [],
  entryComponents: [],
  providers: [
    {
       provide: WiServiceContribution,
       useClass: EventPublisherActivityContributionHandler
     }
  ],
  bootstrap: []
})

export default class EventPublisherActivityModule {

}