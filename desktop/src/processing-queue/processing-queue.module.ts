import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { MatDialogModule } from "@angular/material/dialog";
import { MatIconModule } from "@angular/material/icon";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";
import { MatTooltipModule } from "@angular/material/tooltip";
import { ButtonModule } from "../button";
import { PipesModule } from "../pipes/pipes.module";
import { SharedUiModule } from "../shared-ui/shared-ui.module";
import { TableModule } from "../table/table.module";
import { ProcessingQueueListComponent } from "./processing-queue-list/processing-queue-list.component";
import { ProcessingQueueRoutingModule } from "./processing-queue-routing.module";

/**
 * The Processing page (fork feature), served under `/processing`: quick scans
 * and email uploads that are queued, running, retrying or failed.
 */
@NgModule({
  declarations: [ProcessingQueueListComponent],
  imports: [
    ButtonModule,
    CommonModule,
    MatDialogModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatTooltipModule,
    PipesModule,
    ProcessingQueueRoutingModule,
    SharedUiModule,
    TableModule,
  ],
})
export class ProcessingQueueModule {}
