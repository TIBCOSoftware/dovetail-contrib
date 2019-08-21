import {NgModule} from "@angular/core"
import {HttpModule} from "@angular/http";
import {WiServiceContribution} from "wi-studio/app/contrib/wi-contrib"
import {R3FlowReceiverTriggerHandler} from "./trigger"

@NgModule({
    imports: [
        HttpModule
    ],
    providers: [
        {
            provide: WiServiceContribution,
            useClass: R3FlowReceiverTriggerHandler
        }
    ]
})
export default class R3FlowReceiverTriggerModule {

}
