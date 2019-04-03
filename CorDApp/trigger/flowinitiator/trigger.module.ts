import {NgModule} from "@angular/core"
import {HttpModule} from "@angular/http";
import {WiServiceContribution} from "wi-studio/app/contrib/wi-contrib"
import {R3FlowInitiatorTriggerHandler} from "./trigger"

@NgModule({
    imports: [
        HttpModule
    ],
    providers: [
        {
            provide: WiServiceContribution,
            useClass: R3FlowInitiatorTriggerHandler
        }
    ]
})
export default class R3FlowInitiatorTriggerModule {

}
